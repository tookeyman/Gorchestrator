package netcore

type Client struct {
	socket *socketWorker
}

func GetClientInstance() Client {
	cli := Client{socket: GetSocketInstance()}

	cli.socket.run()
	return cli
}

func (cli *Client) Run() {
	//do nothing
}
