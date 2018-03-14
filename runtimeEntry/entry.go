package main

import (
	"eqbchook2/netcore"
)

func main() {
	a := netcore.GetClientInstance()
	a.Run()
}
