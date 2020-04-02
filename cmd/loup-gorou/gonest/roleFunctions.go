package gonest

import (
	"math/rand"
	"time"
)

//ShuffleRoles func
func ShuffleRoles(playerCount int) []Role {
	roles := make([]Role, 0, playerCount)
	werewolves, villagers := GetRoleDistribution(playerCount)
	for i := 0; i < werewolves; i += 1 {
		roles = append(roles, Role_WEREWOLFROLE)
	}
	for i := 0; i < villagers; i += 1 {
		roles = append(roles, Role_HUMANROLE)
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
