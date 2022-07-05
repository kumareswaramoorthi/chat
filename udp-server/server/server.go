package server

import (
	"fmt"
	"log"
	"net"

	"github.com/kumareswaramoorthi/chat/udp-server/listner"
	"github.com/kumareswaramoorthi/chat/udp-server/message"
	"github.com/kumareswaramoorthi/chat/udp-server/redisclient"
)

// package server provides an interface for UDP server functions.

type Server interface {
	GetNewServer() (*server, error)
	SendMessage()
	HandleMessage()
}

type server struct {
	Conn        *net.UDPConn
	Messages    chan message.Message
	Clients     map[string]Client
	Run         bool
	listner     listner.NetListner
	redisClient redisclient.RedisClient
}

type Client struct {
	userName string
	userAddr *net.UDPAddr
}

func NewServer(listner listner.NetListner, redisClient redisclient.RedisClient) Server {
	return &server{
		listner:     listner,
		redisClient: redisClient,
	}
}

// GetNewServer returns server object.
func (srv *server) GetNewServer() (*server, error) {
	var s server
	var err error
	s.Messages = make(chan message.Message)
	s.Clients = make(map[string]Client)
	s.Run = true
	s.redisClient = srv.redisClient

	s.Conn, err = srv.listner.GetConn()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// SendMessage broadcasts messages to all the active clients.
func (srv *server) SendMessage() {
	for srv.Run {
		msg := <-srv.Messages
		sendMsg := fmt.Sprintf("BROADCAST|time=%s|user-name=%s|message-id=%s|message=%s ", msg.TimeStamp, msg.UserName, msg.MessageID, msg.Content)
		for _, c := range srv.Clients {
			if c.userName != msg.UserName {
				_, err := srv.Conn.WriteToUDP([]byte(sendMsg), c.userAddr)
				if err != nil {
					log.Fatal("unable to send messages")
				}
			}
		}
	}
}

// HandleMessage handles the message from clients.
func (srv *server) HandleMessage() {
	for srv.Run {
		var buf [512]byte
		n, addr, err := srv.Conn.ReadFromUDP(buf[0:])
		if err != nil {
			log.Fatal("unable to read messages")
		}

		m := message.ParseMessage(string(buf[0:n]))

		//check for the action of the incoming message
		switch m.Action {

		// if the client wants to register to the server, the action will be CONNECT.
		// server addes the client address to it's list of clients.
		case "CONNECT":
			if _, ok := srv.Clients[m.UserName]; !ok {
				newClient := Client{
					userAddr: addr,
					userName: m.UserName,
				}
				srv.Clients[m.UserName] = newClient
				history, err := srv.redisClient.RetrieveHistory()
				if err != nil {
					log.Fatal("unable to retrieve history")
				}
				_, err = srv.Conn.WriteToUDP([]byte(history), newClient.userAddr)
				if err != nil {
					log.Fatal("unable to send messages to client")
				}
			}

		// if the client wants to disconnect, the client sends DISCONNECT message.
		// the client address is removed from the server's client lists
		// if all the clients are disconnected then the db is flushed.
		case "DISCONNECT":
			if _, ok := srv.Clients[m.UserName]; ok {
				delete(srv.Clients, m.UserName)
				if len(srv.Clients) == 0 {
					srv.redisClient.Flush()
				}
			}

		// if the client want to delete a particular message, client sends a DELETE action along with the message id.
		// the message is deleted from the redis
		// a notification is sent to all the active clients with the message-id to delete the message.
		case "DELETE":
			err := srv.redisClient.DeleteKey(m)
			if err != nil {
				log.Println("error deleting message")
			}
			for _, c := range srv.Clients {
				msg := fmt.Sprintf("DELETE|%s", m.MessageID)
				_, err := srv.Conn.WriteToUDP([]byte(msg), c.userAddr)
				if err != nil {
					log.Println("unable to send messages to client")
				}
			}

		// if the client wants to deliver message to other clients, it sends BROADCAST action.
		case "BROADCAST":
			err := srv.redisClient.SaveToDB(m)
			if err != nil {
				log.Println("unable to save message to db")
			}
			srv.Messages <- m
		}
	}
}
