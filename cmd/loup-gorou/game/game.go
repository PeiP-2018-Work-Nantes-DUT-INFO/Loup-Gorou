package game

import (
	"errors"
	"sort"

	"github.com/looplab/fsm"
)

type Player struct {
	Name  string
	g     *Game
	Alive bool
}

type CurrentPlayer struct {
	PlayerProps *Player
	Role        string
}

type Game struct {
	Players      map[string]Player
	AlivePlayers map[string]Player
	FSM          *fsm.FSM
	Me           CurrentPlayer
	votes        map[string]string
}

func NewGame(me string, playersNames []string, callbacks fsm.Callbacks) *Game {
	players := make(map[string]Player, len(playersNames))
	for _, name := range playersNames {
		players[name] = Player{
			Name:  name,
			Alive: true,
		}
	}
	g := &Game{
		Players:      players,
		AlivePlayers: make(map[string]Player, len(players)),
		Me: CurrentPlayer{
			PlayerProps: &Player{
				Name: me,
			},
		},
	}
	for _, player := range players {
		g.AlivePlayers[player.Name] = player
		player.g = g
	}
	g.Me.PlayerProps.g = g

	g.FSM = fsm.NewFSM(

		INITIAL_STATE,
		fsm.Events{
			{Name: INITIAL_STATE, Src: []string{}, Dst: WEREWOLF_PLAYING_STATE},
			{Name: WEREWOLF_PLAYING_STATE, Src: []string{VOTE_STATE}, Dst: VOTE_STATE},
			{Name: VOTE_STATE, Src: []string{WEREWOLF_PLAYING_STATE}, Dst: WEREWOLF_PLAYING_STATE},
			{Name: INITIAL_STATE, Src: []string{VOTE_STATE, WEREWOLF_PLAYING_STATE}, Dst: INITIAL_STATE},
		},
		callbacks,
	)
	return g
}

var Fsm *fsm.FSM

func (*Game) KillPlayer(name string) {

}

/*type voteSum struct {
	pseudo string
	sum    int
}

func (g *Game) MostVotedPerson() (p Player, err error) {
	if !g.FSM.Is(VOTE_STATE) || !g.FSM.Is(WEREWOLF_PLAYING_STATE) {
		err = fsm.InvalidEventError{Event: "vote", State: g.FSM.Current()}
		return
	}
	sVotePlay := make(map[string]int)
	for _, target := range g.votes {
		sVotePlay[target] += 1
	}
	var values []voteSum

	for pseudo, value := range sVotePlay {
		values = append(values, voteSum{
			pseudo, value,
		})
	}

	sort.Slice(values[:], func(i, j int) bool {
		return values[i].sum > values[j].sum
	})

	if len(values) < 1 || (len(values) > 1 && values[0].sum == values[1].sum) {

	} else {

	}

	return
}*/
func (g *Game) MostVotedPerson() (p Player, err error) {
	if !g.FSM.Is(VOTE_STATE) || !g.FSM.Is(WEREWOLF_PLAYING_STATE) {
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

func (*Player) Vote(player Player) error {
	if !player.g.FSM.Is(VOTE_STATE) {
		return errors.New("Current state is not vote state")
	}
	player.g.votes[player.Name] = player.Name
	return nil
}

func (*Player) VoteWerewolf(player Player) error {
	if !player.g.FSM.Is(WEREWOLF_PLAYING_STATE) {
		return errors.New("Current state is not werewolf playing state")
	}
	player.g.votes[player.Name] = player.Name
	return nil
}
