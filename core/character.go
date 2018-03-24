package core

import ()

type Character struct {
	*Paperdoll
	cli *Client
}

func GetCharacterInstance(a *Paperdoll, cli *Client) *Character {
	return &Character{a, cli}
}

//Returns the evaluated string of whatever you push to the character
func (char *Character) Query(s string) string {
	char.cli.submitAsyncQuery(char, s, &char.asyncChannel)
	response := <-char.asyncChannel
	return response
}

func (char *Character) DoCommand(callback func()) {
	callback()
}
