package core

import (
	"time"
)

type Character struct {
	*PaperDoll
	cli   *Client
	props *settings
}

func GetCharacterInstance(a *PaperDoll, cli *Client) *Character {
	props := ReadSettingsFile(a.Name)
	char := Character{a, cli, props}
	if props.isDefault {
		sitAndDoNothing := func() {
			char.cli.SendCommandToCharacter(char.Name, "/bca ${Me}s setting file is not configured. "+
				"Please edit ${Me}.props and restart the orchestrator")
			time.Sleep(tick)
		}
		go func() {
			for i := 0; i < 10; i++ {
				char.DoCommand(sitAndDoNothing)
			}
		}()
	}
	return &char
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
