package main

import (
	"bufio"
	"fmt"
	"log"
	"loupgorou/cmd/loup-gorou/gonest"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tidwall/evio"
	"google.golang.org/protobuf/proto"
)

var (
	rightSet     bool = false
	right        net.Conn
	lanIPAddress string
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

	loops, err := strconv.Atoi(os.Getenv("GOROU_EVIO_NUM_LOOPS"))
	if err != nil {
		panic(err.Error())
	}

	var events evio.Events

	switch os.Getenv("GOROU_EVIO_NUM_LOOPS") {
	case "RANDOM":
		events.LoadBalance = evio.Random
	case "ROUND-ROBIN":
		events.LoadBalance = evio.RoundRobin
	case "LEAST-CONNECTIONS":
		events.LoadBalance = evio.LeastConnections
	}
	events.NumLoops = loops
	events.Serving = func(srv evio.Server) (action evio.Action) {
		log.Printf("server started: %s", os.Getenv("GOROU_BIND_ADDRESS"))
		return
	}
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		out = in
		return
	}
	events.Opened = func(ec evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		fmt.Printf("opened: %v\n", ec.RemoteAddr())
		if !rightSet {
			action = evio.Close
			return
		}

		message := &gonest.Event{
			MessageType: gonest.MessageType_ITSHIM,
			Body: &gonest.Event_ItsHimMessage{
				ItsHimMessage: &gonest.ItsHimMessage{},
			},
			IpAddress: lanIPAddress,
		}

		if right.RemoteAddr().String() == "127.0.0.1" {
			message.GetItsHimMessage().RightNodeIpAddress = lanIPAddress
		} else {
			message.GetItsHimMessage().RightNodeIpAddress = right.RemoteAddr().String()
		}

		out, err := proto.Marshal(message)
		if err != nil {
			log.Fatalln("Failed to encode message:", err)
		}

		return
		//ec.SetContext(&conn{})
	}
	events.Closed = func(ec evio.Conn, err error) (action evio.Action) {
		fmt.Printf("closed: %v\n", ec.RemoteAddr())
		return
	}
	if err := evio.Serve(events, os.Getenv("GOROU_BIND_ADDRESS")); err != nil {
		panic(err.Error())
	}

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
		fmt.Println("trying to connect to " + text)
		right, err = net.Dial("tcp", text)
		if err != nil {
			log.Println(err)
		} else {
			rightSet = true
			fmt.Println("connection success")
		}
	}
}
