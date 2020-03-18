package main

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"log"
	"loupgorou/cmd/loup-gorou/gonest"
	"net"
	"os"
	"strings"

	"github.com/firstrow/tcp_server"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"
)

var (
	rightSet     bool = false
	right        net.Conn
	lanIPAddress string
	isLocalhost  bool = true
)

func getIPAdress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func init() {
	lanIPAddress = getIPAdress()
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	server := tcp_server.New(fmt.Sprintf("%s:%s",
		os.Getenv("GOROU_BIND_ADDRESS"),
		os.Getenv("GOROU_BIND_PORT")))

	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		fmt.Println("Message:", message)
		// new message received
	})
	server.OnNewClient(func(c *tcp_server.Client) {
		fmt.Printf("opened: %v\n", c.Conn().RemoteAddr().String())
		message := &gonest.Event{
			MessageType: gonest.MessageType_ITSHIM,
			Body: &gonest.Event_ItsHimMessage{
				ItsHimMessage: &gonest.ItsHimMessage{},
			},
			IpAddress: lanIPAddress,
		}
		if right == nil || strings.Split(right.RemoteAddr().String(), ":")[0] == "127.0.0.1" {
			message.GetItsHimMessage().RightNodeIpAddress = lanIPAddress + ":" + os.Getenv("GOROU_BIND_PORT")
		} else {
			message.GetItsHimMessage().RightNodeIpAddress = right.RemoteAddr().String()
		}

		if !isLocalhost {
			right.Close()
			fmt.Println("Connecting")
			var err error
			host, _, _ := net.SplitHostPort(c.Conn().RemoteAddr().String())
			right, err = net.Dial(c.Conn().RemoteAddr().Network(), "["+host+"]:5000")
			if err != nil {
				panic(err.Error())
			}
		} else {
			isLocalhost = false
		}
		out, err := proto.Marshal(message)
		if err != nil {
			log.Fatalln("Failed to encode message:", err)
		}
		err = c.Send(b64.StdEncoding.EncodeToString(out) + "\n")
		if err != nil {
			panic(err.Error())
		}
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		fmt.Printf("closed: %v\n", c.Conn().RemoteAddr().String())
		// connection with client lost
	})
	go startConnection()
	server.Listen()
}

func startConnection() {
	for !rightSet {
		fmt.Print("Enter ip adress:port (enter for localhost): ")

		//recuperation de l'adresse IP
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		text = strings.TrimSuffix(text, "\r")

		if text == "" {
			text = os.Getenv("GAROU_DEFAULT_CONNECT_ADDRESS")
		}
		fmt.Println("Trying to connect to " + text)
		var err error
		rightSet = true
		right, err = net.Dial("tcp", text)
		if err != nil {
			rightSet = false
			log.Println(err)
		} else {
			fmt.Println("Connection success")
		}
	}

	//New client part
	if right.RemoteAddr().String() != "[::1]:"+os.Getenv("GOROU_BIND_PORT") {
		isLocalhost = false
		reader := bufio.NewReader(right)
		str, err := reader.ReadString('\n')
		if err != nil {
			panic(err.Error())
		}
		buf, err := b64.StdEncoding.DecodeString(str)
		if err != nil {
			panic(err.Error())
		}
		newMsg := &gonest.Event{}
		if err := proto.Unmarshal(buf, newMsg); err != nil {
			panic(err.Error())
		} else {
			fmt.Println(newMsg.String())
		}
		right.Close()
		right, err = net.Dial("tcp", newMsg.GetItsHimMessage().GetRightNodeIpAddress())
		if err != nil {
			panic(err.Error())
		}
		writer := bufio.NewWriter(right)
		writer.WriteString("Test")
		writer.Flush()
	}
}
