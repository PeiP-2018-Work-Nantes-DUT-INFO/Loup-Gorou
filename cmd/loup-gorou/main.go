package main

import (
	"bufio"
	"fmt"
	"log"
	"loupgorou/cmd/loup-gorou/frame"
	"loupgorou/cmd/loup-gorou/gonest"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/firstrow/tcp_server"
	"github.com/joho/godotenv"
)

var (
	rightSet          bool = false
	right             net.Conn
	lanIPAddress      string
	preparationPhase              = true  //création de l'anneau, jusqu'a ce que la partie commence
	completationPhase             = false //lorsque le nombre minimal de personne pour une partie est present dans le réseau, lance un décompte de 20sec avant le debut de la partie
	rightMutex        *sync.Mutex = &sync.Mutex{}
	listIPAddess                  = make([]string, 10)
	minPlayer                     = 3
	ackList           []bool
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
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	if len(os.Args) > 1 {
		os.Setenv("GOROU_BIND_PORT", os.Args[1])
	}
	lanIPAddress = getIPAdress() + ":" + os.Getenv("GOROU_BIND_PORT")
}

func main() {
	server := tcp_server.New(fmt.Sprintf("%s:%s",
		os.Getenv("GOROU_BIND_ADDRESS"),
		os.Getenv("GOROU_BIND_PORT")))

	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		// new message received
		event, err := frame.DecodeEventB64(strings.TrimSuffix(message, "\n"))
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println("Received message", event.String())
		}

		if event.GetSource() != lanIPAddress {
			//rightMutex.Lock()
			switch event.GetMessageType() {
			case gonest.MessageType_HELLO:
				hello(c, event)
				/* tests */
				/* event := gonest.AckMessageFactory(lanIPAddress)
				writer := bufio.NewWriter(right)
				encoded, _ := frame.EncodeEventB64(event)
				_, err := writer.WriteString(encoded)
				if err != nil {
					panic(err.Error())
				}
				err = writer.Flush()
				if err != nil {
					panic(err.Error())
				} */

			case gonest.MessageType_ACK:
				ackHandler()
			case gonest.MessageType_IPLIST:
				ipListHandler(c, event)
			}
			//rightMutex.Unlock()
		} else {
			switch event.GetMessageType() {
			case gonest.MessageType_IPLIST:
				listIPAddess = event.GetIpListMessage().IpAdress
				fmt.Println("actual ip list :", listIPAddess)
				if !preparationPhase {
					ackList = make([]bool, len(listIPAddess))
					sendACK()
				}
			}
		}

	})
	server.OnNewClient(func(c *tcp_server.Client) {
		fmt.Printf("opened: %v\n", c.Conn().RemoteAddr().String())
		if !rightSet {
			c.Close()
		}
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		fmt.Printf("closed: %v\n", c.Conn().RemoteAddr().String())
		// connection with client lost
	})
	go startConnection()
	server.Listen()

}

func hello(c *tcp_server.Client, event *gonest.Event) {
	message := gonest.ItsHimMessageFactory(lanIPAddress)
	if right == nil || strings.Split(right.RemoteAddr().String(), ":")[0] == "127.0.0.1" {
		message.GetItsHimMessage().RightNodeIpAddress = lanIPAddress + ":" + os.Getenv("GOROU_BIND_PORT")
	} else {
		message.GetItsHimMessage().RightNodeIpAddress = right.RemoteAddr().String()
	}
	rightMutex.Lock()
	err := right.Close()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connecting")
	right, err = net.Dial(c.Conn().RemoteAddr().Network(), event.GetSource())
	if err != nil {
		panic(err.Error())
	}
	rightMutex.Unlock()

	encoded, err := frame.EncodeEventB64(message)
	if err != nil {
		log.Fatalln("Failed to encode message:", err)
	}
	err = c.Send(encoded)
	if err != nil {
		panic(err.Error())
	}
}

func ipListHandler(c *tcp_server.Client, event *gonest.Event) {
	event.GetIpListMessage().IpAdress = append(event.GetIpListMessage().IpAdress, lanIPAddress)
	//fmt.Println(event.String())
	encoded, err := frame.EncodeEventB64(event)
	if err != nil {
		panic(err.Error())
	}
	writer := bufio.NewWriter(right)
	_, err = writer.WriteString(encoded)
	if err != nil {
		panic(err.Error())
	}
	err = writer.Flush()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("flush done")
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
		rightMutex.Lock()

		event := gonest.HelloMessageFactory(lanIPAddress)
		writer := bufio.NewWriter(right)
		encoded, _ := frame.EncodeEventB64(event)
		_, err := writer.WriteString(encoded)
		if err != nil {
			panic(err.Error())
		}
		writer.Flush()
		reader := bufio.NewReader(right)
		fmt.Println("Waiting ...")
		str, err := reader.ReadString('\n')
		if err != nil {
			panic(err.Error())
		}
		newMsg, err := frame.DecodeEventB64(str)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println(newMsg.String())
		}
		fmt.Println("Now setting right peer to", newMsg.GetItsHimMessage().GetRightNodeIpAddress())
		err = right.Close()
		if err != nil {
			panic(err.Error())
		}
		right.Close()
		right = nil
		right, err = net.Dial("tcp", newMsg.GetItsHimMessage().GetRightNodeIpAddress())
		if err != nil {
			panic(err.Error())
		}
		rightMutex.Unlock()
	}
	go timerFunction()
}

func timerFunction() {
	tickerIPList := time.NewTicker(5 * time.Second)
	tickerCountMembers := time.NewTicker(1 * time.Second)
	timerEndPreparation := time.NewTimer(1 * time.Hour)

	for {
		select {
		case <-tickerIPList.C:
			if preparationPhase {
				sendIPList()
			}
		case <-tickerCountMembers.C:
			if len(listIPAddess) > minPlayer && preparationPhase && !completationPhase {
				timerEndPreparation = time.NewTimer(20 * time.Second)
				completationPhase = true // debut de la phase de completation
			}
		case <-timerEndPreparation.C:
			if completationPhase {
				completationPhase = false // fin de la pahse de completation
				preparationPhase = false  // fin de la phase de preparation
				sendIPList()
			}
		}
	}
}

func sendACK() {
	event := gonest.AckMessageFactory(lanIPAddress)
	encoded, _ := frame.EncodeEventB64(event)
	rightMutex.Lock()
	writer := bufio.NewWriter(right)
	rightMutex.Unlock()
	_, err := writer.WriteString(encoded)
	if err != nil {
		panic(err.Error())
	}
	writer.Flush()
}

func sendIPList() {
	event := gonest.IpListMessageFactory(lanIPAddress, lanIPAddress)
	encoded, _ := frame.EncodeEventB64(event)
	rightMutex.Lock()
	writer := bufio.NewWriter(right)
	rightMutex.Unlock()
	_, err := writer.WriteString(encoded)
	if err != nil {
		panic(err.Error())
	}
	writer.Flush()
}

func ackHandler() {
	if ackList == nil {
		ackList = make([]bool, len(listIPAddess))
	}
}
