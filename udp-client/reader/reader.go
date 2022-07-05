package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kumareswaramoorthi/chat/udp-client/uuid"
)

type StdInReader interface {
	GetuserName() (string, error)
	GetMessage() (string, error)
	GetMessageID() (string, error)
}

type stdInReader struct {
	uuid uuid.UuidGen
}

func NewStdInReader(uuid uuid.UuidGen) StdInReader {
	return &stdInReader{
		uuid: uuid,
	}
}

// GetuserName reads the user name from the StdIn
func (s *stdInReader) GetuserName() (string, error) {
	var userName string
	read := bufio.NewReader(os.Stdin)
	fmt.Println("enter your user name: ")
	if rawText, err := read.ReadString('\n'); true {
		if err != nil {
			return "", err
		}
		userUUID, _ := s.uuid.GetUUID()
		if err != nil {
			return "", err
		}
		userName = fmt.Sprintf("%s-%d", strings.Trim(rawText, " \n"), userUUID)
		fmt.Printf("here is your user id : %s \n", userName)
	}
	return userName, nil
}

// GetMessage reads the message from StdIn
func (s *stdInReader) GetMessage() (string, error) {
	var msg string
	read := bufio.NewReader(os.Stdin)
	if rawText, err := read.ReadString('\n'); true {
		if err != nil {
			return "", err
		}
		msg = strings.Trim(rawText, " \n")
	}
	return msg, nil
}

// GetMessageID reads the message-id from the StdIn
func (s *stdInReader) GetMessageID() (string, error) {
	var msgID string
	read := bufio.NewReader(os.Stdin)
	fmt.Println("enter message id: ")
	if rawText, err := read.ReadString('\n'); true {
		if err != nil {
			return "", err
		}
		msgID = strings.Trim(rawText, " \n")

	}
	return msgID, nil
}
