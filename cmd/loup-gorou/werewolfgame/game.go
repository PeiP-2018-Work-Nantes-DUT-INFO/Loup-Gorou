package werewolfgame

import (
	"errors"
	"loupgorou/cmd/loup-gorou/gonest"
	"math/rand"
	"sort"
	"time"

	"github.com/looplab/fsm"
)

// This channel is used to know when a new vote is made, an restart the timeout.
var voteChannel chan bool = make(chan bool)

type Player struct {
	Name  string
	g     *Game
	Alive bool
}

type CurrentPlayer struct {
	PlayerProps *Player
	Role        gonest.Role
}

type DeadPlayer struct {
	PlayerProps *Player
	Role        gonest.Role
}

//CanVote
//test if the current player is able to vote given the actual state
func (c *CurrentPlayer) CanVote() bool {
	return (c.PlayerProps.g.FSM.Is(NIGHT_WEREWOLF_PLAYING_STATE) && c.Role == gonest.Role_WEREWOLFROLE) || c.PlayerProps.g.FSM.Is(DAY_VOTE_STATE)
}

func (c *CurrentPlayer) AmIDead() bool {
	if _, ok := c.PlayerProps.g.DeadPlayers[c.PlayerProps.Name]; ok != false {
		return true
	} else {
		return false
	}
}

type Game struct {
	// Contains all the current play of the game
	Players map[string]*Player
	// Contains all the Alive players
	AlivePlayers map[string]*Player
	// Finite State Machine
	FSM *fsm.FSM
	// Represent the current player
	Me CurrentPlayer
	// Contains all the votes
	votes map[string]string
	// List all the death that happened during the night
	MorningDeaths map[string]*Player

	DeadPlayers map[string]*DeadPlayer
}

// Create a new instance of a game.
func NewGame(me CurrentPlayer, playersNames []string, callbacks fsm.Callbacks) *Game {
	players := make(map[string]*Player, len(playersNames))
	for _, name := range playersNames {
		players[name] = &Player{
			Name:  name,
			Alive: true,
		}
	}
	g := &Game{
		Players:       players,
		AlivePlayers:  make(map[string]*Player, len(players)),
		Me:            me,
		votes:         make(map[string]string, len(players)),
		DeadPlayers:   make(map[string]*DeadPlayer, len(players)),
		MorningDeaths: make(map[string]*Player, 5),
	}
	for _, player := range players {
		g.AlivePlayers[player.Name] = player
		player.g = g
	}
	g.Me.PlayerProps.g = g

	// TODO.
	/*callbacks["enter_"+NIGHT_WEREWOLF_PLAYING_STATE] = func(e *fsm.Event) {
		args := make([]interface{}, 1)
		args[0] = NORMAL_END_OF_VOTE
		e.Args = args
		//go g.voteCallback(e)
		callbacks["enter_"+NIGHT_WEREWOLF_PLAYING_STATE](e)
	} */

	/*callbacks["enter_"+DAY_VOTE_STATE] = func(e *fsm.Event) {
		args := make([]interface{}, 1)
		args[0] = NORMAL_END_OF_VOTE
		e.Args = args
		//go g.voteCallback(e)
		callbacks["enter_"+DAY_VOTE_STATE](e)
	}*/

	g.FSM = fsm.NewFSM(

		INITIAL_STATE,
		fsm.Events{
			{Name: START_TRANSITION, Src: []string{INITIAL_STATE}, Dst: NIGHT_WEREWOLF_PLAYING_STATE},
			{Name: WEREWOLF_VOTE_END_TRANSITION, Src: []string{NIGHT_WEREWOLF_PLAYING_STATE}, Dst: DAY_VOTE_STATE},
			{Name: END_OF_DAY_TRANSITION, Src: []string{DAY_VOTE_STATE}, Dst: NIGHT_WEREWOLF_PLAYING_STATE},
			{Name: ALLWEREWOLF_KILLED_DURING_VOTE_TRANSITION, Src: []string{DAY_VOTE_STATE}, Dst: ENDOFGAME_STATE},
			{Name: ALLHUMAN_KILLED_DURING_VOTE_TRANSITION, Src: []string{DAY_VOTE_STATE}, Dst: ENDOFGAME_STATE},
			{Name: ALLWEREWOLF_KILLED_DURING_NIGHT_TRANSITION, Src: []string{DAY_VOTE_STATE}, Dst: ENDOFGAME_STATE},
			{Name: ALLHUMAN_KILLED_DURING_NIGHT_TRANSITION, Src: []string{DAY_VOTE_STATE}, Dst: ENDOFGAME_STATE},
		},
		callbacks,
	)
	return g
}

var Fsm *fsm.FSM

//KillPlayer
//delete the player from the AlivePlayer table
// Used to kill a player when the vote of the day ended
func (g *Game) KillPlayer(player *Player) error {
	_, ok := g.AlivePlayers[player.Name]
	if !ok {
		return errors.New("The player is alread dead")
	}
	delete(g.AlivePlayers, player.Name)
	return nil
}

// Confirm the death of a player
func (g *Game) ConfirmDeath(playerName string, role gonest.Role) {
	g.DeadPlayers[playerName] = &DeadPlayer{
		g.Players[playerName], role,
	}
}

func (g *Game) DoesEveryoneVoted() bool {
	if g.FSM.Is(DAY_VOTE_STATE) {
		return len(g.AlivePlayers) == len(g.votes)
	} else if g.FSM.Is(NIGHT_WEREWOLF_PLAYING_STATE) {
		return len(g.votes) == g.NumberOfAliveWerewolfs()
	} else {
		return false
	}
}

//EatPlayer
//move the player from the alive table to the morning death table
func (g *Game) EatPlayer(player *Player) error {
	_, ok := g.AlivePlayers[player.Name]
	if !ok {
		return errors.New("The player is alread dead")
	}
	g.MorningDeaths[player.Name] = g.AlivePlayers[player.Name]
	return nil
}

//DecideVote
//apply the vote decision (if it's a werewolf vote or a human vote)
func (g *Game) DecideVote() (p *Player, err error) {
	p, err = g.MostVotedPerson()
	if err != nil {
		return
	}
	if g.FSM.Is(NIGHT_WEREWOLF_PLAYING_STATE) {
		err = g.EatPlayer(p)
	} else {
		err = g.KillPlayer(p)
	}
	g.ClearVote()
	return
}

// TODO
/*func (g *Game) voteCallback(e *fsm.Event) {
	var transition string
	if g.FSM.Is(NIGHT_WEREWOLF_PLAYING_STATE) {
		transition = WEREWOLF_VOTE_END_TRANSITION
	} else {
		transition = END_OF_DAY_TRANSITION
	}
	var end bool = false
	for !end {
		select {
		case <-time.After(time.Second * 30):
			//args := make([]interface{}, 1)
			//args[0] = NORMAL_END_OF_VOTE
			//e.Args = args

			_ = g.FSM.Event(transition)
			end = true
		case <-time.After(time.Second * 30):
			if g.isCurrentVoteMajorityAbsolute() {
				_ = g.FSM.Event(transition)
				end = true
			}
		}
	}
}*/

func (g *Game) IsCurrentVoteMajorityAbsolute() bool {
	sVotePlay := make(map[string]int)
	for _, target := range g.votes {
		sVotePlay[target] += 1
	}
	return len(sVotePlay) == 1
}

//MostVotedPerson
//get the most voted person. Handle equality case by comparing the ip adress
func (g *Game) MostVotedPerson() (p *Player, err error) {
	if !g.FSM.Is(DAY_VOTE_STATE) && !g.FSM.Is(NIGHT_WEREWOLF_PLAYING_STATE) {
		err = fsm.InvalidEventError{Event: "vote", State: g.FSM.Current()}
		return
	}

	sVotePlay := make(map[string]int)
	for _, target := range g.votes {
		sVotePlay[target] += 1
	}

	sVotePlay_2 := make(map[int][]string)
	var values []int
	for key, value := range sVotePlay {
		sVotePlay_2[value] = append(sVotePlay_2[value], key)
		values = append(values, value)
	}

	_ = sort.Reverse(sort.IntSlice(values))

	var playersToChoose []string
	if len(values) < 1 || len(sVotePlay_2[values[0]]) > 1 {
		if len(values) < 1 {
			keys := make([]string, 0, len(g.AlivePlayers))
			for k := range g.AlivePlayers {
				keys = append(keys, k)
			}
			playersToChoose = keys
		} else {
			playersToChoose = sVotePlay_2[values[0]]
		}
		sort.SliceStable(playersToChoose, func(i, j int) bool {
			return playersToChoose[i] < playersToChoose[j]
		})
		p = g.AlivePlayers[playersToChoose[0]]
	} else {
		p = g.AlivePlayers[sVotePlay_2[values[0]][0]]
	}
	return
}

//clearVote
//reset the vote table
func (g *Game) ClearVote() {
	g.votes = make(map[string]string)
}

//NumberOfVotes
//get the number of vote
func (g *Game) NumberOfVotes() int {
	return len(g.votes)
}

func (g *Game) NumberOfAliveWerewolfs() int {
	werewolfes, _ := GetRoleDistribution(len(g.Players))
	for _, deadPlayer := range g.DeadPlayers {
		if deadPlayer.Role == gonest.Role_WEREWOLFROLE {
			werewolfes -= 1
		}
	}
	return werewolfes
}

func (g *Game) WerewolfWon() bool {
	if len(g.AlivePlayers) == g.NumberOfAliveWerewolfs() {
		return true
	} else {
		return false
	}
}

func (g *Game) HumansWon() bool {
	return g.NumberOfAliveWerewolfs() == 0
}

func (g *Game) GetMorningDeaths() []*Player {
	deadPlayers := make([]*Player, 0, len(g.MorningDeaths))
	for _, player := range g.MorningDeaths {
		deadPlayers = append(deadPlayers, player)
		delete(g.AlivePlayers, player.Name)
	}
	g.MorningDeaths = make(map[string]*Player)
	return deadPlayers
}

//Vote
//Send a vote message
func (c *Player) Vote(name string) error {
	player, ok := c.g.AlivePlayers[name]
	if !ok {
		return errors.New("Player is dead")
	}
	if !player.g.FSM.Is(DAY_VOTE_STATE) && !player.g.FSM.Is(NIGHT_WEREWOLF_PLAYING_STATE) {
		return errors.New("Current state is not a vote state")
	}
	// voteChannel <- true
	player.g.votes[c.Name] = player.Name
	return nil
}

func ShuffleRoles(playerCount int) []gonest.Role {
	roles := make([]gonest.Role, 0, playerCount)
	werewolves, villagers := GetRoleDistribution(playerCount)
	for i := 0; i < werewolves; i += 1 {
		roles = append(roles, gonest.Role_WEREWOLFROLE)
	}
	for i := 0; i < villagers; i += 1 {
		roles = append(roles, gonest.Role_HUMANROLE)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(roles), func(i, j int) { roles[i], roles[j] = roles[j], roles[i] })
	return roles
}

func GetRoleDistribution(playerCount int) (int, int) {
	if playerCount <= 8 {
		return 1, playerCount - 1
	} else if playerCount <= 11 {
		return 2, playerCount - 2
	} else if playerCount <= 17 {
		return 3, playerCount - 3
	} else {
		return 4, playerCount - 4
	}
}
