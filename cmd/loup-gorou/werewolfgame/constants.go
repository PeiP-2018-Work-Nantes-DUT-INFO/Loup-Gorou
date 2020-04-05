package werewolfgame

//definition of the differents state of the game
const (
	INITIAL_STATE                = "INITIAL_STATE"  //starting state
	ENDOFGAME_STATE              = "END_STATE"      //ending state
	NIGHT_WEREWOLF_PLAYING_STATE = "WEREWOLF_STATE" //werewolf turn state. This state allow the werewolf to make votes
	DAY_VOTE_STATE               = "VOTE_STATE"     //day state. This state allow the villagers to makes votes
)

//definition of differents transition
const (
	START_TRANSITION                           = "NEWGAME_TRANSITION"
	WEREWOLF_VOTE_END_TRANSITION               = "WEREWORLF_END_TRANDITION"
	ALLWEREWOLF_KILLED_DURING_VOTE_TRANSITION  = "ALLWEREWOLF_KILLED_DURING_VOTE_TRANSITION"
	ALLHUMAN_KILLED_DURING_VOTE_TRANSITION     = "ALLHUMAN_KILLED_DURING_VOTE_TRANSITION"
	ALLWEREWOLF_KILLED_DURING_NIGHT_TRANSITION = "ALLWEREWOLF_KILLED_DURING_NIGHT_TRANSITION"
	ALLHUMAN_KILLED_DURING_NIGHT_TRANSITION    = "ALLHUMAN_KILLED_DURING_VOTE_TRANSITION"
	END_OF_DAY_TRANSITION                      = "END_OF_DAY_TRANSITION"
)

const (
	NORMAL_END_OF_VOTE = 0
	TIMEOUT_OF_VOTE    = 1
)
