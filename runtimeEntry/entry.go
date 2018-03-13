package main

import(
	"eqbchook2/netcore"
)

func main() {
	a:= netcore.GetSocketInstance()
	a.Write(netcore.Login.String())
	a.Test()
}
