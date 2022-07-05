package main

import (
	"fmt"
	"log"

	"github.com/kumareswaramoorthi/chat/udp-client/client"
	"github.com/kumareswaramoorthi/chat/udp-client/dialer"
	"github.com/kumareswaramoorthi/chat/udp-client/reader"
	"github.com/kumareswaramoorthi/chat/udp-client/uuid"
)

func main() {

	uuid := uuid.NewUuidGen()
	dialer := dialer.NewNetDialer()
	reader := reader.NewStdInReader(uuid)
	client := client.NewChatClient(uuid, reader, dialer)

	newUdpClient, err := client.GetNewChatClient()
	if err != nil {
		log.Fatal("internal server error")
	}

	err = newUdpClient.Connect()
	if err != nil {
		log.Fatal("unable to connect to server")
	}

	go newUdpClient.PrintMessage()
	go newUdpClient.ReceiveMessage()
	go newUdpClient.SendMessage()
	newUdpClient.GetMessage()

	fmt.Println("chat disconnected")
}
