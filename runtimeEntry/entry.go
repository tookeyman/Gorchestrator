package main

import (
	"eqbchook2/core"
	"fmt"
	"time"
)

func main() {
	client := core.GetClientInstance()
	//client.Run()
	fmt.Println("Waiting for ten seconds for reasons")
	time.Sleep(1 * time.Second)
	client.Test()
}
