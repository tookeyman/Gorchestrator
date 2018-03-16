package main

import (
	"eqbchook2/netcore"
	"fmt"
	"time"
)

func main() {
	client := netcore.GetClientInstance()
	client.Run()
	fmt.Println("Waiting for ten seconds for reasons")
	time.Sleep(10 * time.Second)

}
