package client

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/kumareswaramoorthi/chat/udp-client/dialer"
	"github.com/kumareswaramoorthi/chat/udp-client/reader"
	"github.com/kumareswaramoorthi/chat/udp-client/uuid"
)

// package client provides an interface for UDP client functions.

type ChatClient interface {
	GetNewChatClient() (*chatClient, error)
	Connect() error
	PrintMessage()
	SendMessage()
	GetMessage()
}

type chatClient struct {
	conn            *net.UDPConn
	userName        string
	sendMessages    chan string
	receiveMessages chan string
	run             bool
	uuid            uuid.UuidGen
	stdInReader     reader.StdInReader
	dialer          dialer.NetDialer
}

func NewChatClient(uuid uuid.UuidGen, stdInReader reader.StdInReader, dialer dialer.NetDialer) ChatClient {
	return &chatClient{
		uuid:        uuid,
		stdInReader: stdInReader,
		dialer:      dialer,
	}
}

// GetNewChatClient provides a new client object
func (client *chatClient) GetNewChatClient() (*chatClient, error) {
	var c chatClient
	var err error
	c.sendMessages = make(chan string)
	c.receiveMessages = make(chan string)
	c.userName, err = client.stdInReader.GetuserName()

	if err != nil {
		return nil, err
	}

	c.conn, err = client.dialer.GetConn()
	if err != nil {
		return nil, err
	}

	c.uuid = uuid.NewUuidGen()
	c.dialer = dialer.NewNetDialer()
	c.stdInReader = reader.NewStdInReader(c.uuid)

	c.run = true
	return &c, nil

}

// Connect sends a connect action message to the chat server
func (client *chatClient) Connect() error {
	str := fmt.Sprintf("CONNECT|%s", client.userName)
	_, err := client.conn.Write([]byte(str))
	if err != nil {
		return err
	}
	return nil
}

// PrintMessage prints the message.
func (client *chatClient) PrintMessage() {
	for client.run {
		msg := <-client.receiveMessages
		fmt.Println(msg)
	}
}

// GetMessage reads the message from user through StdIn.
func (client *chatClient) GetMessage() {
	fmt.Println("You have entered chat room. You can start chatting now ...")
	for client.run {
		msg, err := client.stdInReader.GetMessage()
		if err != nil {
			log.Fatal("unable to read message")
		}

		switch msg {
		// if the user want to quit the chat
		case ":quit":
			str := fmt.Sprintf("DISCONNECT|%s", client.userName)
			_, err := client.conn.Write([]byte(str))
			if err != nil {
				log.Fatal("unable to send message")
			}
			return
			// if user wants to delete a chat message.
		case ":delete":
			msgID, err := client.stdInReader.GetMessageID()
			if err != nil {
				log.Fatal("unable to read message")
			}
			str := fmt.Sprintf("DELETE|%s|%s", client.userName, msgID)
			_, err = client.conn.Write([]byte(str))
			if err != nil {
				log.Fatal("unable to send message")
			}
			fmt.Println("your message has been deleted")

			// by default the message is treated as chat message and sent to the chat server.
		default:
			client.sendMessages <- msg
		}
	}
}

// SendMessage sends the message to the chat server.
func (client *chatClient) SendMessage() {
	for client.run {
		msg := <-client.sendMessages
		uuid, err := client.uuid.GetUUID()
		if err != nil {
			log.Fatal("unable to generate uuid for message")
		}
		str := fmt.Sprintf("BROADCAST|%s|%d|%s", client.userName, uuid, msg)
		_, err = client.conn.Write([]byte(str))
		if err != nil {
			log.Fatal("unable to send message")
		}

		// notification is show to the user that the message has been sent and it's unique id.
		fmt.Printf("[notification : time=%s|message sent|message-id=%d|message=%s] \n", time.Now().Format("2006-01-02 3:4:5 pm"), uuid, msg)
	}
}

// ReceiveMessage receives the message from the chat server
func (client *chatClient) ReceiveMessage() {
	var buf [4096]byte
	for client.run {
		n, err := client.conn.Read(buf[0:])
		if err != nil {
			log.Fatal("unable to read message")
		}
		msg := strings.Split(string(buf[0:n]), "|")
		switch msg[0] {
		// if the action is to delete a message, the a notification is shown to the user.
		case "DELETE":
			client.receiveMessages <- fmt.Sprintf("[notification : delete message-id=%s]", msg[1])
			// if the action was to broadcast, then the actual message is shown to the user.
		case "BROADCAST":
			client.receiveMessages <- fmt.Sprintf("%s|%s", strings.Join(msg[1:3], "|"), strings.Join(msg[4:], "|"))

		default:
			client.receiveMessages <- strings.Join(msg, "|")
		}
	}
}
