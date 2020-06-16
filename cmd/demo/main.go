package main

import (
	"fmt"
	"log"

	"github.com/maxmcd/cave"
	. "github.com/maxmcd/cave/dom"
)

type State struct {
	value string
}

func (state State) Render() Node {
	ss := StringState("value", state.value, state.onInputChange)

	return Div(
		Div(Input(OnChange(ss.Callback))),
		Div(ss.Value),
	)
}

func (state State) onInputChange(value string) {
	state.value = value
}

func main() {
	server := cave.New(func() cave.Renderable {
		return State{value: "hello"}
	})
	addr := ":8080"
	fmt.Println("listening on", addr)
	log.Fatal(server.ListenAndServe(addr))
}
