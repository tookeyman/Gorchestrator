package core

import (
	"fmt"
	"regexp"
	"sync"
	"time"
)

const (
	tick = 6 * time.Second
)

type Client struct {
	socket           *socketWorker
	asyncQue         map[string]*chan string
	asyncQueMutex    *sync.Mutex
	characters       map[string]*Character
	requestWatermark int
}

//test method
func (cli *Client) Test() {

	for cli.socket.running {

	}
}

//creates a client, initializes it and returns its pointer
func GetClientInstance() *Client {

	cli := Client{
		asyncQue:         make(map[string]*chan string, 50), /*i dunno, seems good*/
		asyncQueMutex:    &sync.Mutex{},
		characters:       make(map[string]*Character),
		requestWatermark: 0,
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
	cli.requestWatermark++
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
		cha := GetPaperdollInstance(groups[0], groups[1])
		cli.characters[groups[0]] = GetCharacterInstance(cha, cli)
	} else {
		cha := cli.characters[groups[0]]
		cha.UpdatePaperdoll(groups[1])
	}
}

//sends out the async part of the async/await. THIS ISN'T THREAD SAFE, IT CAN ONLY HANDLE ABOUT 250 REQ/S
func (cli *Client) submitAsyncQuery(char *Character, request string, stringChannel *chan string) {
	payload := fmt.Sprintf("//bct Orchestrator [ASYNC]%s", request)
	cli.asyncQueMutex.Lock()
	cli.asyncQue[char.Name] = stringChannel
	cli.asyncQueMutex.Unlock()
	go func() {
		cli.SendCommandToCharacter(char.Name, payload)
	}()
}

//is the await part of the character query async/await
func (cli *Client) handleAsyncResponse(response [2]string) {
	cli.asyncQueMutex.Lock()
	stringHandle := cli.asyncQue[response[0]]
	cli.asyncQueMutex.Unlock()
	if *stringHandle != nil {
		*stringHandle <- response[1]
	}
	delete(cli.asyncQue, response[0])
}
