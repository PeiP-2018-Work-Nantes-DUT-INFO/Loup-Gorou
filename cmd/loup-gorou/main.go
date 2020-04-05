package main

import (
	"bufio"
	"fmt"
	"loupgorou/cmd/loup-gorou/frame"
	"loupgorou/cmd/loup-gorou/gonest"
	"loupgorou/cmd/loup-gorou/secondstimer"
	"loupgorou/cmd/loup-gorou/werewolfgame"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/firstrow/tcp_server"
	"github.com/joho/godotenv"
	"github.com/looplab/fsm"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Those constants below are used in the finite state machine
const (
	PREPARATION_STATE      = "PREPARATION_STATE"      //preparation state designates the moment when the ring network does not contain enough members
	COMPLETATION_STATE     = "COMPLETATION_STATE"     //completation state designates the moment when they are enough peoples in the network ring.
	LEADERELECTION_STATE   = "LEADERELECTION_STATE"   //leaderelection state designates the moment when the leader is elected.
	ROLEDISTRIBUTION_STATE = "ROLEDISTRIBUTION_STATE" //roledistribution state designates the moment when the leader is giving the role to the other player.
	GAME_STATE             = "GAME_STATE"             //game state designates the moment when the game is running.
)

const (
	CONNECTED_TRANSITION       = "CONNECTED_TRANSITION"       //connected transition is the transition uses to go to the COMPLETATION_STATE from the PREPARATION_STATE.
	RING_COMPLETED_TRANSITION  = "RING_COMPLETED_TRANSITION"  //ringcompleted transition is the transition uses to go to the LEADERELECTION_STATE from the COMPLETATION_STATE.
	LEADER_ELECTED_TRANSITION  = "LEADER_ELECTED_TRANSITION"  //leaderelected transition is the transition uses to go to the ROLEDISTRIBUTION_STATE from the LEADERELECTION_STATE.
	ROLE_DISTRIBUED_TRANSITION = "ROLE_DISTRIBUED_TRANSITION" //roledistributed transition is the transition uses to go to the GAME_STATE from the ROLEDISTRIBUTION_STATE.

	TIME_BEFORE_TIMEOUT_MATCHMAKING_SEC = 30
)

var (
	gamemaster          *log.Entry
	rightSet            bool                       = false              //define if the ring is create
	right               net.Conn                                        //define the connection to our right neightboor int the ring network.
	lanIPAddress        string                                          //contains our ip address with the tcp server port
	rightMutex          *sync.Mutex                = &sync.Mutex{}      //mutex to protect the right connection
	listIPAddress                                  = make([]string, 10) //contains every ip address of the players in the ring.
	minPlayer                                      = 3                  // the value is set at 2 for test
	ackMap              map[string]int                                  //in key we have the ip addresses of the players and the value is the result is if they acquit the state.
	timerEndPreparation *secondstimer.SecondsTimer                      //timer start when we have the minimum number of player in the ring, when the timer is finished, the game start. (the time will be redifined in the future).
	leader              string                                          //contains the leader ip address
	gameInstance        *werewolfgame.Game                              //gameInstance define the current game situation, that struct contains the finals states machine of the game.
	fsmConnection       *fsm.FSM                                        //finals states machine describe the connection phase.
)

//getIPAdress
//return the IP adress of the first interface
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

//promptUser
//get the input of the user since the begining of the game
//and take action in consequence
func promptUser() {
	for {
		fmt.Print(">")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		text = strings.TrimSuffix(text, "\r")
		if gameInstance.Me.AmIDead() {
			gamemaster.Error("You are dead")
			return
		}
		if text != "" {
			if text[0] == '!' {
				if len(text) > 6 && text[1:5] == "vote" {
					target := text[6:]

					//check the ip address validity
					if _, ok := gameInstance.AlivePlayers[target]; !ok {
						gamemaster.Error("Invalid target")
						continue
					}
					if gameInstance.Me.CanVote() {
						sendVote(target)
					} else {
						gamemaster.Error("You are not allowed to vote right now")
					}

				}
			} else {
				if gameInstance.FSM.Is(werewolfgame.DAY_VOTE_STATE) {
					sendChatMessage(text)
				}
			}
		}
	}
}

//init
//initialisation of the different Callback
func init() {
	log.SetFormatter(&prefixed.TextFormatter{})
	log.SetLevel(log.TraceLevel)
	gamemaster = log.WithField("prefix", "GAMEMASTER")
	// load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	// if we have an argument, then set the GOROU_BIND_PORT env variable to his value.
	// Allow use to start the application by specifing directly a port number
	if len(os.Args) > 1 {
		os.Setenv("GOROU_BIND_PORT", os.Args[1])
	}
	//get the local IP address of the first interface that has one. Should work on Linux, hazardous on Windows.
	lanIPAddress = getIPAdress() + ":" + os.Getenv("GOROU_BIND_PORT")

	// Initialize the CONNNECTION FSM
	// This state machine is used to describe sequentially the different "phases" of the ring establishment.
	fsmConnection = fsm.NewFSM(PREPARATION_STATE, []fsm.EventDesc{
		fsm.EventDesc{Name: CONNECTED_TRANSITION, Src: []string{PREPARATION_STATE}, Dst: COMPLETATION_STATE},
		fsm.EventDesc{Name: RING_COMPLETED_TRANSITION, Src: []string{COMPLETATION_STATE}, Dst: LEADERELECTION_STATE},
		fsm.EventDesc{Name: LEADER_ELECTED_TRANSITION, Src: []string{LEADERELECTION_STATE}, Dst: ROLEDISTRIBUTION_STATE},
		fsm.EventDesc{Name: ROLE_DISTRIBUED_TRANSITION, Src: []string{ROLEDISTRIBUTION_STATE}, Dst: GAME_STATE},
	},
		//initialisation of the different callback of the state machine
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				log.WithField("prefix", "CONNECTION FSM").Debug("Entering state ", e.Dst) // will be printed each time we enter a state
			},
			"leave_state": func(e *fsm.Event) {
				if e.FSM.Is(COMPLETATION_STATE) || e.FSM.Is(LEADERELECTION_STATE) || e.FSM.Is(ROLEDISTRIBUTION_STATE) {
					e.Async()
					sendACK()
				}
				log.WithField("prefix", "CONNECTION FSM").Debug("Leaving state ", e.Src) // will be printed each time we leave a state
			},
			"leave_" + PREPARATION_STATE: func(e *fsm.Event) {
				log.WithField("prefix", "CONNECTION FSM").Info("Finished preparation phase !") // will be printed when the preparation state is finished.
			},
			"enter_" + COMPLETATION_STATE: func(e *fsm.Event) {
				go timerFunction()
				log.WithField("prefix", "CONNECTION FSM").Info("Beginning completation phase")
			},
		},
	)
}

func main() {
	//Initialisation of the tcp server.
	server := tcp_server.New(fmt.Sprintf("%s:%s",
		os.Getenv("GOROU_BIND_ADDRESS"),
		os.Getenv("GOROU_BIND_PORT")))

	//handler activate when a new message is received on the tcp server.
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		// new message received
		event, err := frame.DecodeEventB64(strings.TrimSuffix(message, "\n"))
		if err != nil {
			log.Error("Could not decode event ", err)
			return
		} else {
			log.Trace(event.String())
		}

		if event.GetSource() != lanIPAddress {
			broadcastInProgress(c, message, event)
		} else {
			broadcastDone(message, event)
		}

	})
	//handler activate when a new client open a connection to the tcp server.
	server.OnNewClient(func(c *tcp_server.Client) {
		log.Infof("opened: %v", c.Conn().RemoteAddr().String())
		if !rightSet {
			c.Close()
		}
	})
	//handler executed when a close the connection with the server.
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
		log.Infof("closed: %v", c.Conn().RemoteAddr().String())
		if fsmConnection.Is(GAME_STATE) || fsmConnection.Is(ROLEDISTRIBUTION_STATE) {
			if fsmConnection.Is(GAME_STATE) && gameInstance.FSM.Is(werewolfgame.ENDOFGAME_STATE) {
				return
			}
			log.Fatal("Unexcepted connection close")
		}
	})
	go startConnection() // start the prompt
	//start listen and handling connections.
	server.Listen() // blocking function
}

// broadcastInProgress is used when a message is sent to this client but is not from us. (we are not the message source).
// So all those handlers are supposed to relay the messages
func broadcastInProgress(c *tcp_server.Client, message string, event *gonest.Event) {
	switch event.GetMessageType() {
	case gonest.MessageType_HELLO:
		helloHandler(c, event)
	case gonest.MessageType_ACK:
		ackHandler(event)
	case gonest.MessageType_IPLIST:
		ipListHandler(event)
	case gonest.MessageType_CHAT:
		chatHandler(event)
	case gonest.MessageType_ROLEDISTRIBUTION:
		roleHandler(event)
	case gonest.MessageType_LEADERELECTION:
		leaderElectionHandler(event)
	case gonest.MessageType_VOTE:
		voteHandler(event)
	case gonest.MessageType_DEAD:
		deadHandler(event)
	}
}

// broadcastDone is called when a message is sent to this client and is from us. (we are the message source).
// It basically mean that the message turned all over the ring and came back to us
func broadcastDone(message string, event *gonest.Event) {
	switch event.GetMessageType() {
	case gonest.MessageType_IPLIST:
		listIPAddress = event.GetIpListMessage().IpAdress
		if len(listIPAddress) >= minPlayer && fsmConnection.Is(PREPARATION_STATE) {
			_ = fsmConnection.Event(CONNECTED_TRANSITION)
		}
	case gonest.MessageType_ROLEDISTRIBUTION:
		if fsmConnection.Is(ROLEDISTRIBUTION_STATE) {
			role := event.GetRoleDistributionMessage().GetRole()
			if event.GetRoleDistributionMessage().Target == lanIPAddress {
				initGame(role)
			}

		}
	// Handler mixé, utilisé à la fois lorsqu'un broadcast est en progress, et est terminé
	case gonest.MessageType_ACK:
		ackHandler(event)
	case gonest.MessageType_VOTE:
		voteHandler(event)
	case gonest.MessageType_CHAT:
		chatHandler(event)
	case gonest.MessageType_DEAD:
		deadHandler(event)
	}
}

//function use to enter in the ring network or to create it.
func startConnection() {
	for !rightSet {
		fmt.Print("Enter ip adress:port (enter for localhost): ")

		//recuperation de l'adresse IP
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		//text := ""
		text = strings.TrimSuffix(text, "\n")
		text = strings.TrimSuffix(text, "\r")

		if text == "" {
			text = os.Getenv("GAROU_DEFAULT_CONNECT_ADDRESS")
		}
		log.Infof("Trying to connect to %s", text)
		var err error
		rightSet = true
		right, err = net.Dial("tcp", text)
		if err != nil {
			rightSet = false
			log.Error(err)
		} else {
			log.Info("Connection success")
		}
	}

	//New client partd()
	if right.RemoteAddr().String() != "[::1]:"+os.Getenv("GOROU_BIND_PORT") && right.RemoteAddr().String() != "127.0.0.1:"+os.Getenv("GOROU_BIND_PORT") {
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
		log.Info("Waiting ...")
		str, err := reader.ReadString('\n')
		if err != nil {
			panic(err.Error())
		}
		newMsg, err := frame.DecodeEventB64(str)
		if err != nil {
			panic(err.Error())
		} else {
			log.Trace(newMsg.String())
		}
		log.Debugf("Now setting right peer to %s", newMsg.GetItsHimMessage().GetRightNodeIpAddress())
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
	sendIPList()
}

func timerFunction() {
	timerEndPreparation = secondstimer.NewSecondsTimer(TIME_BEFORE_TIMEOUT_MATCHMAKING_SEC * time.Second)
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-timerEndPreparation.Timer.C:
			if fsmConnection.Is(COMPLETATION_STATE) {
				log.Warn("GAME BEGINS")
				_ = fsmConnection.Event(RING_COMPLETED_TRANSITION)
				return
			}
		case <-ticker.C:
			log.Infof("Time remaining before match begins: %.0fs", timerEndPreparation.TimeRemaining().Seconds())
		}

	}
}

func initGame(role gonest.Role) {
	gamemaster.Warn("You have the role ", role.String())
	gameInstance = werewolfgame.NewGame(werewolfgame.NewCurrentPlayer(lanIPAddress, role),
		listIPAddress,
		fsm.Callbacks{
			"leave_state": func(e *fsm.Event) {
				if e.Src != werewolfgame.INITIAL_STATE {
					sendACK()
					e.Async()
				}
				log.WithField("prefix", "GAME FSM").Debug("Leaving state from ", e.Src)
			},
			"enter_state": func(e *fsm.Event) {
				log.WithField("prefix", "GAME FSM").Debug("Entering state ", e.Dst)
			},
			"enter_" + werewolfgame.NIGHT_WEREWOLF_PLAYING_STATE: func(e *fsm.Event) {
				ended := checkVictoryStateHandler(e)
				if !ended {
					gamemaster.Info("The night comes on Thiercelieux ...")
					gamemaster.Info("The werewolves wake up")
					if gameInstance.Me.Role == gonest.Role_WEREWOLFROLE {
						gamemaster.Info("You are a werewolf. Vote with !vote [ipaddress]")
					}
				}
			},
			"leave_" + werewolfgame.NIGHT_WEREWOLF_PLAYING_STATE: func(e *fsm.Event) {
				gamemaster.Info("Morning comes on the village. A ray of sunlight light up the bell tower.")
				deadPlayers := gameInstance.GetMorningDeaths()
				if len(deadPlayers) > 0 {
					gamemaster.Info("People got killed this night :")
					for _, player := range deadPlayers {
						gamemaster.Println("\t\t-", player.Name)
						if player.Name == lanIPAddress {
							sendDead(gameInstance.Me.Role, gonest.Reason_NORMAL)
						}
					}
				} else {
					gamemaster.Println("No one died during this night.")
				}
			},
			"enter_" + werewolfgame.DAY_VOTE_STATE: func(e *fsm.Event) {
				checkVictoryStateHandler(e)
			},
		})
	_ = fsmConnection.Event(ROLE_DISTRIBUED_TRANSITION)
}

func getLeader() string {
	sort.Strings(listIPAddress)
	return listIPAddress[0]
}

func initAckMap() {
	ackMap = make(map[string]int, len(listIPAddress))
	for _, value := range listIPAddress {
		ackMap[value] = 0
	}
}

func isEverybodyOk() (result bool) {
	result = true
	for _, value := range ackMap {
		result = result && value > 0
	}
	if result {
		for key := range ackMap {
			ackMap[key] -= 1
		}
	}
	return
}

func checkVictoryStateHandler(e *fsm.Event) bool {
	funcToexec := func() {
		time.Sleep(time.Second * 3)
		os.Exit(0)
	}
	if gameInstance.WerewolfWon() {
		gamemaster.Info("Werewolfs won")
		go funcToexec()
		gameInstance.FSM.SetState(werewolfgame.ENDOFGAME_STATE)
		return true

	} else if gameInstance.HumansWon() {
		gamemaster.Info("Humans won")
		go funcToexec()
		gameInstance.FSM.SetState(werewolfgame.ENDOFGAME_STATE)
		return true
	}
	return false
}

func giveRoles() {
	roles := werewolfgame.ShuffleRoles(len(listIPAddress))
	for index, player := range listIPAddress {
		log.Debug(roles[index], "send to", player)
		go sendRole(player, roles[index])
	}
}

// MESSAGE HANDLING
func helloHandler(c *tcp_server.Client, event *gonest.Event) {
	message := gonest.ItsHimMessageFactory(lanIPAddress)
	if right == nil || strings.Split(right.RemoteAddr().String(), ":")[0] == "127.0.0.1" {
		message.GetItsHimMessage().RightNodeIpAddress = lanIPAddress
	} else {
		message.GetItsHimMessage().RightNodeIpAddress = right.RemoteAddr().String()
	}
	rightMutex.Lock()
	err := right.Close()
	if err != nil {
		panic(err.Error())
	}
	log.Info("Connecting")
	right, err = net.Dial(c.Conn().RemoteAddr().Network(), event.GetSource())
	if err != nil {
		panic(err.Error())
	}
	rightMutex.Unlock()

	encoded, err := frame.EncodeEventB64(message)
	if err != nil {
		log.Fatal("Failed to encode message")
	}
	err = c.Send(encoded)
	if err != nil {
		panic(err.Error())
	}
}

func deadHandler(event *gonest.Event) {
	if event.GetSource() != lanIPAddress {
		go eventPropagator(event, right)
	}
	deadMessage := event.GetDeadMessage()
	gameInstance.ConfirmDeath(event.GetSource(), deadMessage.GetRole())
	gamemaster.Warnf("%s was %s", event.GetSource(), deadMessage.GetRole())

}

func ipListHandler(event *gonest.Event) {
	table := event.GetIpListMessage().GetIpAdress()
	event.GetIpListMessage().IpAdress = append(table, lanIPAddress)
	go eventPropagator(event, right)
	isPresent := false
	for _, ip := range listIPAddress {
		isPresent = isPresent || (ip == event.GetSource())
	}
	if !isPresent {
		listIPAddress = append(listIPAddress, event.GetSource())
	}
	if len(listIPAddress) >= minPlayer && fsmConnection.Is(PREPARATION_STATE) {
		_ = fsmConnection.Event(CONNECTED_TRANSITION)
	} else if fsmConnection.Is(COMPLETATION_STATE) {
		timerEndPreparation.Reset(TIME_BEFORE_TIMEOUT_MATCHMAKING_SEC * time.Second)
	}
}

func roleHandler(event *gonest.Event) {
	go eventPropagator(event, right)
	roleMessage := event.GetRoleDistributionMessage()
	if roleMessage.Target == lanIPAddress && fsmConnection.Is(ROLEDISTRIBUTION_STATE) {
		initGame(roleMessage.GetRole())
	}
}

func voteHandler(event *gonest.Event) {
	if event.GetSource() != lanIPAddress {
		go eventPropagator(event, right)
	}
	vote := event.GetVoteMessage()
	err := gameInstance.AlivePlayers[event.GetSource()].Vote(vote.GetTarget())
	if err != nil {
		log.Error(err)
	}
	if gameInstance.FSM.Is(werewolfgame.DAY_VOTE_STATE) {
		gamemaster.Info(event.GetSource(), " voted for ", vote.GetTarget()+".")
	}
	if gameInstance.DoesEveryoneVoted() {
		player, err := gameInstance.DecideVote()
		if err != nil {
			panic(err.Error())
		}
		if gameInstance.FSM.Is(werewolfgame.NIGHT_WEREWOLF_PLAYING_STATE) {
			_ = gameInstance.FSM.Event(werewolfgame.WEREWOLF_VOTE_END_TRANSITION)
		} else {
			if gameInstance.Me.PlayerProps.Name == player.Name {
				gamemaster.Warn("You are dead")
				sendDead(gameInstance.Me.Role, gonest.Reason_NORMAL)
			} else {
				gamemaster.Info(player.Name, "is dead.")
			}
			_ = gameInstance.FSM.Event(werewolfgame.END_OF_DAY_TRANSITION)
		}
	}
}

func ackHandler(event *gonest.Event) {
	if event.Source != lanIPAddress {
		go eventPropagator(event, right)
	}
	if ackMap == nil {
		initAckMap()
	}
	source := event.GetSource()
	ackMap[source] += 1
	if isEverybodyOk() {
		log.Debug("Everybody is okay. ACK")
		if fsmConnection.Is(COMPLETATION_STATE) || fsmConnection.Is(LEADERELECTION_STATE) || fsmConnection.Is(ROLEDISTRIBUTION_STATE) {
			fsmConnection.Transition()
		}
		if fsmConnection.Is(LEADERELECTION_STATE) {
			leader = getLeader()
			if leader == lanIPAddress {
				log.Debug("I'm the leader ", leader)
				event := gonest.LeaderElectionMessageFactory(lanIPAddress, leader)
				eventPropagator(event, right)
			}
			_ = fsmConnection.Event(LEADER_ELECTED_TRANSITION)
		} else if fsmConnection.Is(ROLEDISTRIBUTION_STATE) {
			leader = getLeader()
			if leader == lanIPAddress {
				giveRoles()
			}
		} else if fsmConnection.Is(GAME_STATE) {
			switch gameInstance.FSM.Current() {
			case werewolfgame.INITIAL_STATE:
				//roleDistributionPhase = false
				gamemaster.Warn("STARTING GAME")
				gamemaster.Info("Player in game are:")
				for _, value := range listIPAddress {
					fmt.Println("\t\t-", value)
				}
				_ = gameInstance.FSM.Event(werewolfgame.START_TRANSITION)
				go promptUser()
				/*
					} else if fsmConnection.Is(NIGHT_WEREWOLF_PLAYING_STATE) {
						if (isAllWerewolfDead()) {
							_ = gameInstance.FSM.Event(werewolfgame.ALLWEREWOLF_KILLED_DURING_NIGHT_TRANSITION)
						} else if (isAllVillagerDead()) {
							_ = gameInstance.FSM.Event(werewolfgame.ALLHUMAN_KILLED_DURING_NIGHT_TRANSITION)
						} else {
							_ = gameInstance.FSM.Event(werewolfgame.WEREWORLF_END_TRANDITION)
						}

					} else if fsmConnection.Is(DAY_VOTE_STATE) {
						if (isAllWerewolfDead()) {
							_ = gameInstance.FSM.Event(werewolfgame.ALLWEREWOLF_KILLED_DURING_VOTE)
						} else if (isAllVillagerDead()) {
							_ = gameInstance.FSM.Event(werewolfgame.ALLHUMAN_KILLED_DURING_VOTE)
						} else {
							_ = gameInstance.FSM.Event(werewolfgame.END_OF_DAY_TRANSITION)
						}
				*/
			default:
				err := gameInstance.FSM.Transition()
				if err != nil {
					log.Error(err)
				}

			}
		}
	}

}

func chatHandler(event *gonest.Event) {
	if event.GetSource() != lanIPAddress {
		go eventPropagator(event, right)
	}
	contentMessage := event.GetChatMessage().GetContent()
	source := event.GetSource()

	log.WithField("prefix", "CHAT").Println(source, ":", contentMessage)
}

func leaderElectionHandler(event *gonest.Event) {
	leader := event.GetLeaderElectionMessage().GetLeader()
	eventPropagator(event, right)
	/*if !fsmConnection.Is(LEADERELECTION_STATE) {
		fsmConnection.SetState(ROLEDISTRIBUTION_STATE)
		sendACK()
	}*/
	log.Info("Leader elected is ", leader, ", now awaiting from him to get my role !")
}

// MESSAGE CREATION AND SEND
func eventPropagator(event *gonest.Event, right net.Conn) {
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

func sendChatMessage(message string) {
	//if actual state is day state
	event := gonest.ChatMessageFactory(lanIPAddress, message)
	eventPropagator(event, right)
}

func sendACK() {
	if ackMap == nil {
		initAckMap()
	}
	event := gonest.AckMessageFactory(lanIPAddress)
	eventPropagator(event, right)
}

func sendDead(role gonest.Role, reason gonest.Reason) {
	event := gonest.DeadMessageFactory(lanIPAddress, role, reason)
	eventPropagator(event, right)
}

func sendVote(target string) {
	event := gonest.VoteMessageFactory(lanIPAddress, target)
	eventPropagator(event, right)
}

func sendIPList() {
	event := gonest.IpListMessageFactory(lanIPAddress, lanIPAddress)
	eventPropagator(event, right)
}

func sendRole(target string, role gonest.Role) {
	event := gonest.RoleDistributionMessageFactory(lanIPAddress, target, role)
	eventPropagator(event, right)
}
