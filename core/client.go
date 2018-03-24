package core

import (
	"fmt"
	"regexp"
	"time"
)

const (
	tick = 6 * time.Second
)

type Client struct {
	socket     *socketWorker
	asyncQue   map[string]*chan string
	characters map[string]*Character
}

//test method
func (cli *Client) Test() {
	fired := false
	for cli.socket.running {
		if !fired && len(cli.characters) > 0 {
			for key := range cli.characters {
				char := cli.characters[key]
				char.Benchmark()
			}
			fired = true
		}
		//fmt.Println("Waiting for one minute...")
		//for i := 0; i < 10; i++ {
		//for key := range cli.characters {
		//	fmt.Printf("[%s]\t%#v\n", key, *cli.characters[key])
		//}
		//	time.Sleep(tick)
		//}
		//fmt.Println("Minute finished...")
		//cli.socket.ToggleRunning()
		//cli.Disconnect()
	}
}

//creates a client, initializes it and returns its pointer
func GetClientInstance() *Client {
	cli := Client{
		asyncQue:   make(map[string]*chan string, 50), /*i dunno, seems good*/
		characters: make(map[string]*Character),
	}
	sock, err := GetSocketInstance(cli.handleSocketRead)
	if err != nil {
		fmt.Printf("[CLIENT]\tClient instance creation failed: %s\n", err)
	}
	cli.socket = sock
	return &cli
}

//probably poorly named. sends the command string to the character name
func (cli *Client) SendCommandToCharacter(characterName string, command string) {
	packet := fmt.Sprintf("%s%s %s", Tell.String(), characterName, command)
	cli.socket.Write(packet)
}

//this function is called every time we get a string in off of the socket. primary source of logic for how to route info
func (cli *Client) handleSocketRead(message string) {
	netBots := regexp.MustCompile(`^.*NBPKT:(\w+):\[NB]\|(.*)\[NB]\n$`)
	asyncResponse := regexp.MustCompile(`^\[(\w+)] \[ASYNC](.*)\n$`)

	switch message {
	case "\tPING\n":
		cli.socket.Write(Pong.String())
		return
	default:
	}

	if netBots.MatchString(message) {
		matches := netBots.FindStringSubmatch(message)
		netBotsPayload := [2]string{matches[1], matches[2]}
		cli.handleNetbotsPacket(netBotsPayload)
	} else if asyncResponse.MatchString(message) {
		//@blocking this has to block for each character being queried, or else it will dereference a null pointer
		matches := asyncResponse.FindStringSubmatch(message)
		asyncPayload := [2]string{matches[1], matches[2]}
		cli.handleAsyncResponse(asyncPayload)
	} else {
		fmt.Println(message)
	}
}

//sends a disconnect packet to the eqbcserver
func (cli *Client) Disconnect() {
	cli.socket.Write(Disconnect.String())
}

//broadcasts as if you were doing a bca/a
func (cli Client) Broadcast(message string) {
	cli.socket.Write(fmt.Sprintf("%s %s", MsgAll, message))
}

//client method for routing netbot packets
func (cli *Client) handleNetbotsPacket(groups [2]string) {
	if cli.characters[groups[0]] == nil {
		cha := GetActorInstance(groups[0], groups[1])
		cli.characters[groups[0]] = GetCharacterInstance(cha, cli)
	} else {
		cha := cli.characters[groups[0]]
		cha.UpdateActor(groups[1])
	}
}

//sends out the async part of the async/await
func (cli *Client) submitAsyncQuery(char *Character, request string, stringHandle *chan string) {
	payload := fmt.Sprintf("//squelch /bct Orchestrator [ASYNC]%s", request)
	cli.asyncQue[char.Name] = stringHandle
	go func() {
		cli.SendCommandToCharacter(char.Name, payload)
	}()
}

//is the await part of the character query async/await
func (cli *Client) handleAsyncResponse(response [2]string) {
	stringHandle := cli.asyncQue[response[0]]
	*stringHandle <- response[1]
	delete(cli.asyncQue, response[0])
}
