package core

import (
	"fmt"
	"strings"
	"time"
)

type (
	CharacterState uint8
	Command        struct {
		validStates     []CharacterState
		commandCallBack func()
	}
	Character struct {
		*PaperDoll
		cli          *Client
		Props        *settings
		CommandQue   chan *Command
		CurrentState CharacterState
	}
)

const (
	ANY_STATE CharacterState = iota
	REST
	FOLLOW
	COMBAT
)

func GetCharacterInstance(a *PaperDoll, cli *Client) *Character {
	props := ReadSettingsFile(a.Name)
	char := &Character{a, cli, props, make(chan *Command, 50), REST}
	go char.initiateCommandQueue()
	char.SubmitCommand(createPropertiesValidationCommand(char))
	return char
}

func (char *Character) Benchmark() {
	fmt.Println("Starting benchmark for", char.Name)
	count := 0
	totalTime := int64(0)

	for i := 0; i < 1000; i++ {
		start := time.Now()
		char.Query(fmt.Sprintf("req: %d", i))
		dur := time.Since(start).Nanoseconds()
		count += 1
		totalTime += dur
	}

	fmt.Printf("%s finished %d reqs in %dns of total time. %f avg\n",
		char.Name, count, totalTime, float32(totalTime/int64(count)))
}

func createPropertiesValidationCommand(char *Character) *Command {
	sitAndDoNothing := func() {
		char.cli.SendCommandToCharacter(char.Name, "//bcaa My setting file is not configured. "+
			"Please edit ${Me}.props and restart the orchestrator")
		time.Sleep(tick)
	}
	idleCommand := CreateCharacterCommand(sitAndDoNothing, ANY_STATE)

	checkForDefaultProps := func() {
		if char.Props.isDefault {
			for i := 0; i < 10; i++ {
				char.SubmitCommand(idleCommand)
			}
		}
	}
	return CreateCharacterCommand(checkForDefaultProps, ANY_STATE)
}

func (char *Character) SubmitCommand(command *Command) {
	char.CommandQue <- command
}

func (char *Character) initiateCommandQueue() {
	for char.cli.socket.running {
		select {
		case currentCommand := <-char.CommandQue:
			if currentCommand.matchesState(char.CurrentState) {
				currentCommand.commandCallBack()
			}
		}
	}
}

//Returns the evaluated string of whatever you push to the character
func (char *Character) Query(s string) string {
	char.cli.submitAsyncQuery(char, s, &char.asyncChannel)
	response := <-char.asyncChannel
	return response
}

//
func (char *Character) SlashCommand(s string) {
	char.cli.SendCommandToCharacter(char.Name, "/"+s)
}

func (com *Command) matchesState(state CharacterState) bool {
	for i := range com.validStates {
		if com.validStates[i] == state || com.validStates[i] == ANY_STATE {
			return true
		}
	}
	return false
}

func (char *Character) CreateSpellCommand(spell string, target string, validState CharacterState) *Command {
	callBack := func() {
		for char.casting == "c" {
			time.Sleep(100 * time.Millisecond)
		}
		spellCheck := fmt.Sprintf("${Raid.Member[%s]}@${Group.Member[%s]}@${Spawn[pc %s]}", target, target, target)
		queryPacket := strings.Split(char.Query(spellCheck), "@")
		isRaidMember := queryPacket[0] != "NULL"
		isGroupMember := queryPacket[1] != "NULL"
		isPC := queryPacket[3] != "NULL"

		if isGroupMember {
			char.SlashCommand(fmt.Sprintf("/g Casting %s on %s", spell, target))
		} else if isRaidMember {
			char.SlashCommand(fmt.Sprintf("/raid Casting %s on %s", spell, target))
		} else if isPC {
			char.SlashCommand(fmt.Sprintf("/tell %s Casting %s on you", target, spell))
		}

		spellCastCommand := fmt.Sprintf("/casting %s %s", spell, target)
		char.SlashCommand(spellCastCommand)
		time.Sleep(300 * time.Millisecond)
		for char.casting == "c" {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return CreateCharacterCommand(callBack, validState)
}

func CreateCharacterCommand(callBack func(), validStates ...CharacterState) *Command {
	return &Command{validStates: validStates, commandCallBack: callBack}
}
