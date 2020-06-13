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
	return Div(
		Div(Input(OnChange(StateInput("value", state.onInputChange)))),
		Div(StateOutput("value", state.value)),
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
