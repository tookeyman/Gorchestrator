package netcore

import (
	"fmt"
	"regexp"
	//"strconv"
	"eqbchook2/character"
	"time"
)

type Client struct {
	socket     *socketWorker
	asyncQue   map[string]*chan string
	characters map[string]*character.Character
}

func (cli *Client) Test() {

	fmt.Println("Waiting for one minute...")
	time.Sleep(1 * time.Minute)
	fmt.Println("Minute finished...")
	cli.socket.ToggleRunning()
	cli.Disconnect()
}

//creates a client, initializes it and returns its address
func GetClientInstance() *Client {
	cli := Client{
		asyncQue:   make(map[string]*chan string, 50), /*i dunno, seems good*/
		characters: make(map[string]*character.Character),
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
		matches := asyncResponse.FindStringSubmatch(message)
		asyncPayload := [2]string{matches[1], matches[2]}
		cli.handleAsyncResponse(asyncPayload)
	} else {
		fmt.Println(message)
	}
}

func (cli *Client) Disconnect() {
	cli.socket.Write(Disconnect.String())
}

func (cli Client) Broadcast(message string) {
	cli.socket.Write(fmt.Sprintf("%s %s", MsgAll, message))
}

func (cli *Client) handleNetbotsPacket(groups [2]string) {
	if cli.characters[groups[0]] == nil {
		cha := character.GetCharacterInstance(groups[0], groups[1])
		cli.characters[groups[0]] = cha
	} else {
		cha := *cli.characters[groups[0]]
		cha.UpdateCharacter(groups[1])
	}
	fmt.Printf("[NETBOTS]\tRECEIVED:%#v\n", groups)
}

func (cli *Client) submitAsyncQuery(char string, response string, stringHandle *chan string) {
	//@TODO we could probably make this completely nonblocking, but we can't enforce order without a request index, which we could do
	payload := fmt.Sprintf("//bct Orchestrator [ASYNC]%s", response)
	cli.asyncQue[char] = stringHandle
	go func() {
		cli.SendCommandToCharacter(char, payload)
	}()
}

func (cli *Client) handleAsyncResponse(response [2]string) {
	stringHandle := cli.asyncQue[response[0]]
	*stringHandle <- response[1]
	delete(cli.asyncQue, response[0])
}

func (cli *Client) Run() {
	//do nothing
}
