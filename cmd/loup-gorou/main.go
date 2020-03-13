package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/firstrow/tcp_server"
	"github.com/joho/godotenv"
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
		log.Fatal("Error loading .env file")
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
