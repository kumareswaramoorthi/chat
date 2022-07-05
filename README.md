# UDP CHAT  


- Redis backed chat room written in Go language.
- Chat room based on UDP protocol.
- Clients send messages via UDP to Server. Then it is broadcast to all other clients.
- Server pushes message to Redis for temporary history. History is limited to 20  messages. 
- When new client connects to chat server it receives last 20 messages (in correct  order).
- Client may delete any message he/she has previously written (but not  messages from others). 
- When client deletes message it will also be removed from redis and a notification is send to client to delete the message, new clients will see  history without it. 
- When all clients disconnect, the DB is flushed. 


NOTE : Since the client is written in GO and chat messages are printed in terminal, it's difficult to
remove the previous message from user's view, but a notification is sent from the server to delete the particular 
message by its message-id.



Repository Tree:
-----------------

```
├── README.md
├── go.mod
├── go.sum
├── udp-client
│   ├── client
│   │   └── client.go
│   ├── dialer
│   │   └── dialer.go
│   ├── main.go
│   ├── reader
│   │   └── reader.go
│   └── uuid
│       └── uuid.go
└── udp-server
    ├── listner
    │   └── listner.go
    ├── main.go
    ├── message
    │   └── message.go
    ├── redisclient
    │   └── redis.go
    └── server
        └── server.go
```
	
	
Requirements:
-----------------

1. Go 1.17+ 
2. Redis



How to run:
-----------------

1. Start Redis in port 6379
     
      - docker : `$ docker run --name my-redis -p 6379:6379 -d redis`
      - ubuntu : `$ sudo systemctl start redis.service`


2. Clone the repo

	`$ git clone https://github.com/kumareswaramoorthi/chat.git`

3. Navigate to project directory 

	`$ cd chat`

4. Navigate to the UDP Server folder

	`$ cd udp-server`

5. Run the server by the following command 

	`$ go run main.go`

6. Navigate to the udp-client folder

	`$ cd ../udp-client`

7. Start a client to enter into chat 

	`$ go run main.go`
    
