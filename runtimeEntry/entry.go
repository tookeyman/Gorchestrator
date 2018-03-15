package main

import (
	"eqbchook2/netcore"
	"fmt"
	"time"
)

func main() {
	a := netcore.GetClientInstance()
	a.Run()
	fmt.Println("Waiting for ten seconds for reasons")
	time.Sleep(10 * time.Second)
	a.Test()
}
