package werewolfgame

//definition of the differents state of the game
const (
	INITIAL_STATE                = "INITIAL_STATE"           //starting state
	CUPID_PREPARATION_STATE      = "CUPID_PREPARATION_STATE" // Cupid playing state, play first
	ENDOFGAME_STATE              = "END_STATE"               //ending state
	NIGHT_WEREWOLF_PLAYING_STATE = "WEREWOLF_STATE"          //werewolf turn state. This state allow the werewolf to make votes
	DAY_VOTE_STATE               = "VOTE_STATE"              //day state. This state allow the villagers to makes votes
)

//definition of differents transition
const (
	START_TRANSITION                   = "NEWGAME_TRANSITION"
	CUPID_END_TRANSITION               = "CUPID_END_TRANSITION"
	WEREWOLF_VOTE_END_TRANSITION       = "WEREWORLF_END_TRANDITION"
	END_OF_GAME_AFTER_VOTE_TRANSITION  = "END_OF_GAME_AFTER_VOTE_TRANSITION"
	END_OF_GAME_AFTER_NIGHT_TRANSITION = "END_OF_GAME_AFTER_NIGHT_TRANSITION"
	END_OF_DAY_TRANSITION              = "END_OF_DAY_TRANSITION"
)

const (
	NORMAL_END_OF_VOTE = 0
	TIMEOUT_OF_VOTE    = 1
)
