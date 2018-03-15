package netcore

import (
	"fmt"
	"regexp"
	"time"
)

type Client struct {
	socket *socketWorker
}

func (cli *Client) Test() {

	cli.SendCommandToCharacter("Prizzae", "//bct Orchestrator hello")

	fmt.Println("Waiting for one minute...")
	time.Sleep(1 * time.Minute)
	fmt.Println("Minute finished...")
	cli.socket.ToggleRunning()
}

func GetClientInstance() Client {
	cli := Client{}
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
	//@todo: we need message handling for netbot packets, character struct creation, all the good stuff

	netBots := regexp.MustCompile(`^.*NBPKT:(\w+):\[NB]\|(.*)\[NB]\n$`)
	switch message {
	case "\tPING\n":
		cli.socket.Write(Pong.String())
		return
	default:
	}

	if netBots.MatchString(message) {
		//works
	} else {
		fmt.Println(message)
	}

}

func (cli *Client) Run() {
	//do nothing
}
