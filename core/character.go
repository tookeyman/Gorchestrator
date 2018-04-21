package core

import (
	"time"
)

func GetCharacterInstance(a *PaperDoll, cli *Client) *Character {
	props := ReadSettingsFile(a.Name)
	char := &Character{a, cli, props, make(chan Command, 50), REST}
	go char.initiateCommandQueue()
	char.SubmitCommand(createPropertiesValidationCommand(char))
	return char
}

func createPropertiesValidationCommand(char *Character) *Command {
	sitAndDoNothing := func() {
		char.cli.SendCommandToCharacter(char.Name, "//bcaa My setting file is not configured. "+
			"Please edit ${Me}.props and restart the orchestrator")
		time.Sleep(tick)
	}
	idleCommand := CreateCharacterCommand(sitAndDoNothing, ANY_STATE)

	checkForDefaultProps := func() {
		if char.props.isDefault {
			for i := 0; i < 10; i++ {
				char.SubmitCommand(idleCommand)
			}
		}
	}
	return CreateCharacterCommand(checkForDefaultProps, ANY_STATE)
}

func (char *Character) SubmitCommand(command *Command) {
	char.commandQue <- *command
}

func (char *Character) initiateCommandQueue() {
	for char.cli.socket.running || len(char.commandQue) > 0 {
		currentCommand := <-char.commandQue
		if currentCommand.matchesState(char.currentState) {
			currentCommand.commandCallBack()
		}
	}
}

//Returns the evaluated string of whatever you push to the character
func (char *Character) Query(s string) string {
	char.cli.submitAsyncQuery(char, s, &char.asyncChannel)
	response := <-char.asyncChannel
	return response
}

func (com *Command) matchesState(state CharacterState) bool {
	for i := range com.validStates {
		if com.validStates[i] == state || com.validStates[i] == ANY_STATE {
			return true
		}
	}
	return false
}

func CreateCharacterCommand(callBack func(), validStates ...CharacterState) *Command {
	return &Command{validStates: validStates, commandCallBack: callBack}
}

type Command struct {
	validStates     []CharacterState
	commandCallBack func()
}

type Character struct {
	*PaperDoll
	cli          *Client
	props        *settings
	commandQue   chan Command
	currentState CharacterState
}
type CharacterState uint

const (
	ANY_STATE CharacterState = iota
	REST
	FOLLOW
	COMBAT
)
