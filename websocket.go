package cave

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/maxmcd/cave/internal/messages"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"golang.org/x/net/html"
)

type websocketSession struct {
	req  *http.Request
	conn *websocket.Conn
	cave *Cave

	// TODO: make this not a pointer
	components map[string]*liveComponent
}

type liveComponent struct {
	renderer Renderer
	tree     *html.Node
}

func newLiveState(r Renderer) *liveComponent {
	ls := liveComponent{renderer: r}
	var buf bytes.Buffer
	Render(r, &buf)
	ls.tree = mustParse(buf.String())
	// printNode(ls.tree, 0)
	return &ls
}

func (ls *liveComponent) diff() ([]Patch, error) {
	old := ls.tree
	var buf bytes.Buffer
	Render(ls.renderer, &buf)
	ls.tree = mustParse(buf.String())
	// printNode(ls.tree, 0)
	return Diff(old, ls.tree)
}

type EventType string

const (
	EventTypePatch  EventType = "p"
	EventTypeSubmit EventType = "submit"
)

func (wss *websocketSession) sendError(err error) {
	if err == nil {
		return
	}
	msg, _ := messages.NewServerMessage("", string(messages.ServerTypeError), []string{err.Error()})
	_ = wss.conn.WriteJSON(msg)
}

func (wss *websocketSession) handleMessages() {
	for {
		_, msg, err := wss.conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				return
			}
			// TODO
			panic(err)
		}
		var jsonMsg messages.ClientMessage
		_ = json.Unmarshal(msg, &jsonMsg)
		if err := wss.handleClientMessage(jsonMsg); err != nil {
			wss.sendError(err)
			return
		}
	}
}

func (wss *websocketSession) handleClientMessage(msg messages.ClientMessage) error {
	spew.Dump(msg)
	switch msg.EventType() {
	case messages.ClientTypeInit:
		var componentIDs []string
		if err := msg.UnmarshalBody(&componentIDs); err != nil {
			return err
		}
		if wss.components != nil {
			return errors.New("can't initialize the websocket twice")
		}
		wss.components = map[string]*liveComponent{}
		for _, componentID := range componentIDs {
			parts := strings.Split(componentID, "-")
			if len(parts) == 0 {
				return fmt.Errorf("component id %q is not in the correct format", componentID)
			}
			componentName := parts[0]
			rendererFunc, ok := wss.cave.registry[componentName]
			if !ok {
				return fmt.Errorf("component %q is not in the registry", componentName)
			}
			renderer := rendererFunc()
			wss.components[componentID] = newLiveState(renderer)
			onmounter, ok := renderer.(OnMounter)
			if ok {
				onmounter.OnMount(wss.req)
				return wss.reRenderComponent(componentID)
			}
		}
	case messages.ClientTypeSubmit:
		var data []map[string]string
		if err := msg.UnmarshalBody(&data); err != nil {
			return err
		}
		ls, ok := wss.components[msg.ComponentID()]
		if ok {
			fmt.Println(data)
			ls.renderer.(OnSubmiter).OnSubmit(msg.Name(), data[0])
		}
		return wss.reRenderComponent(msg.ComponentID())
	}
	return nil
}

func (wss *websocketSession) reRenderComponent(componentID string) error {
	component, ok := wss.components[componentID]
	if !ok {
		return fmt.Errorf("component with id %q not found", componentID)
	}
	diff, err := component.diff()
	if err != nil {
		return err
	}
	if len(diff) == 0 {
		return nil
	}
	out, _ := messages.NewServerMessage(
		componentID,
		string(messages.ServerTypePatch),
		diff,
	)
	return wss.conn.WriteJSON(out)
}
