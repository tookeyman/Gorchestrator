package core

import (
	"fmt"
	"time"
)

type Character struct {
	*Actor
	cli *Client
}

func GetCharacterInstance(a *Actor, cli *Client) *Character {
	return &Character{a, cli}
}

func (char *Character) Benchmark() {
	start := time.Now()
	iterations := 1000
	for i := 0; i < iterations; i++ {
		str := fmt.Sprintf("%d", i)
		fmt.Println(char.Query(str))
	}
	seconds := time.Since(start).Seconds()
	fmt.Printf("%.2f queries/s in %.2fs\n", float64(iterations)/seconds, seconds)
}

func (char *Character) Query(s string) string {
	char.cli.submitAsyncQuery(char, s, &char.asyncChannel)
	response := <-char.asyncChannel
	return response
}
