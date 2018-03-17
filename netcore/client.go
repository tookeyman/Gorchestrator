package netcore

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Client struct {
	socket   *socketWorker
	asyncQue map[string]*chan string
}

func (cli *Client) Test() {
	stringHandle := make(chan string)

	for i := 0; i < 100; i++ {
		cli.submitAsyncQuery("Prizzae", strconv.Itoa(i), &stringHandle)
		fmt.Println(<-stringHandle)
	}

	fmt.Println("Waiting for one minute...")
	time.Sleep(1 * time.Minute)
	fmt.Println("Minute finished...")
	cli.socket.ToggleRunning()
	cli.Disconnect()
}

func GetClientInstance() Client {
	cli := Client{asyncQue: make(map[string]*chan string, 50 /*i dunno, seems good*/)}
	sock, err := GetSocketInstance(cli.handleSocketRead)
	if err != nil {
		fmt.Printf("[CLIENT]\tClient instance creation failed: %s\n", err)
	}
	cli.socket = sock
	return cli
}

func (cli *Client) SendCommandToCharacter(characterName string, command string) {
	packet := fmt.Sprintf("%s%s %s", Tell.String(), characterName, command)
	cli.socket.Write(packet)
}

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
	//@todo: we need message handling for netbot packets, character struct creation, all the good stuff
	//fmt.Printf("[NETBOTS]\tRECEIVED:%#v\n",groups)
}

func (cli *Client) submitAsyncQuery(char string, response string, stringHandle *chan string) {
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
