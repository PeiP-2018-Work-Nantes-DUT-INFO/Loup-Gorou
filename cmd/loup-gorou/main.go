package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/firstrow/tcp_server"
	"github.com/joho/godotenv"
	"github.com/tidwall/evio"
)

var (
	rightSet bool = false
	right    net.Conn
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

	return "nil"
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	loops, err := strconv.Atoi(os.Getenv("GOROU_EVIO_NUM_LOOPS"))
	if err != nil {
		panic(err.Error())
	}

	var events evio.Events

	switch os.Getenv("GOROU_EVIO_NUM_LOOPS") {
	case "RANDIM":
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
		//ec.SetContext(&conn{})
		return
	}
	events.Closed = func(ec evio.Conn, err error) (action evio.Action) {
		fmt.Printf("closed: %v\n", ec.RemoteAddr())
		return
	}
	if err := evio.Serve(events, os.Getenv("GOROU_BIND_ADDRESS")); err != nil {
		panic(err.Error())
	}

	server := tcp_server.New(os.Getenv("GAROU_BIND_ADDRESS"))

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter ip adress:port (enter for localhost): ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		fmt.Println("Local party set at " + getIPAdress())
	} else {
		right, err = net.Dial("tcp", text)
		if err != nil {
			log.Println(err)
		} else {
			for {
				reader := bufio.NewReader(right)
				message, err := reader.ReadString('\n')
				if err != nil {
					log.Println(err)
				}
				fmt.Print(message)
			}
			rightSet = true
		}
	}

	server.OnNewClient(func(c *tcp_server.Client) {

		if rightSet == false {
			c.Send(getIPAdress() + "\n")
		} else {
			//on envoie l'adresse ip du membre droit
		}

	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		// new message received
		fmt.Println(message)
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
	})

	server.Listen()
}
