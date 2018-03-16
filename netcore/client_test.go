package netcore

func (cli *Client) TestClient(){
	cli.Broadcast("hello")


	cli.Disconnect()
}
