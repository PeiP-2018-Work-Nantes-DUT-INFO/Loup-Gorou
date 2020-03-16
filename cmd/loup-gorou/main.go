package main

import (
	"fmt"
	"log"
	"loupgorou/cmd/loup-gorou/gonest"
	"os"

	"github.com/firstrow/tcp_server"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	msg := &gonest.Event{
		MessageType: gonest.MessageType_CUPID,
		Body: &gonest.Event_CupidMessage{
			CupidMessage: &gonest.CupidMessage{
				IpAddress1: "12.12.12.12",
				IpAddress2: "13.13.13.13",
			},
		},
		IpAddress: "127.0.0.1",
	}

	fmt.Println(msg.String())
	msg.GetCupidMessage().IpAddress1 = "15.15.15.15"
	msg.GetCupidMessage().IpAddress2 = "16.16.16.16"
	msgMarsh, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}

	newMsg := &gonest.Event{}
	if err := proto.Unmarshal([]byte(msgMarsh), newMsg); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	} else {
		fmt.Println(newMsg.String())
	}

	server := tcp_server.New(os.Getenv("GAROU_BIND_ADDRESS"))

	server.OnNewClient(func(c *tcp_server.Client) {
		// new client connected
		// lets send some message
		c.Send("Hello")
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		// new message received
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
	})

	server.Listen()
}
