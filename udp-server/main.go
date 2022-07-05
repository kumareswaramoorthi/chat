package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/kumareswaramoorthi/chat/udp-server/listner"
	"github.com/kumareswaramoorthi/chat/udp-server/redisclient"
	"github.com/kumareswaramoorthi/chat/udp-server/server"
)

func main() {

	redisClient := redisclient.NewRedisClient(context.Background(), redis.NewClient(&redis.Options{Addr: "localhost:6379"}))
	listner := listner.NewNetListner()
	s, err := server.NewServer(listner, redisClient).GetNewServer()
	if err != nil {
		log.Fatal("internal server error")
	}

	go s.SendMessage()
	go s.HandleMessage()

	fmt.Println("chat server started...")

	// Listen for an OS signal to stop the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("Server exiting")
}
