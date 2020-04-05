// Simon SASSI's part

package werewolfgame

import (
	"errors"
	"loupgorou/cmd/loup-gorou/gonest"
)

const (
	CUPID_ACTION_CONTEXT = "CUPID_ACTION_CONTEXT"
	CUPID_IN_LOVE        = "CUPID_IN_LOVE"
)

type CupidAction struct {
	target1 string
	target2 string
}

func (c *CurrentPlayer) InitCupidContext() error {
	if c.Role != gonest.Role_CUPIDROLE {
		return errors.New("I don't have the CUPID role, I can't inititialize my cupid context")
	}
	c.SetContext(CUPID_ACTION_CONTEXT, &CupidAction{})
	return nil
}

func (c *CurrentPlayer) GetCupidContext() (action *CupidAction, err error) {
	temp, ok := c.GetContext(CUPID_ACTION_CONTEXT)
	if !ok {
		err = errors.New("Context not initialized")
		return
	}
	action = temp.(*CupidAction)
	return
}

func (a *CupidAction) MarkArrowAsSent(name string) bool {
	if a.target1 == "" {
		a.target1 = name
		return true
	} else {
		if name != a.target1 {
			a.target2 = name
			return true
		} else {
			return false
		}
	}
}
func (a *CupidAction) HasSentAllArrows() bool {
	return a.target1 != "" && a.target2 != ""
}

func (a *CupidAction) GetArrows() (string, string, error) {
	if !a.HasSentAllArrows() {
		return "", "", errors.New("You didn't send all arrows")
	}
	return a.target1, a.target2, nil
}

func (c *CurrentPlayer) FallInLoveWith(name string) {
	c.SetContext(CUPID_IN_LOVE, name)
}

func (c *CurrentPlayer) IsInLoveWith(name string) bool {
	temp, ok := c.GetContext(CUPID_IN_LOVE)
	if !ok {
		return false
	}
	myLover := temp.(string)
	if myLover == name {
		return true
	} else {
		return false
	}
}
