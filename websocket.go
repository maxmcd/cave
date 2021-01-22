package cave

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/maxmcd/cave/internal/messages"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"golang.org/x/net/html"
)

var (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
)

type websocketSession struct {
	req  *http.Request
	conn *websocket.Conn
	cave *Cave

	// TODO: make this not a pointer
	components map[string]*liveComponent
}

type componentRenderer struct {
	Renderer
	id int
}

type liveComponent struct {
	renderer Renderer
	tree     *html.Node

	// need to check if we've seen the subcomponent before
	// need to take an action from a server and direct it to the right subcomponent
	// need to select an existing index for rendering
	subcomponents map[uintptr]componentRenderer
	idMap         map[int]uintptr
	counter       int
}

func (lc *liveComponent) render(renderer Renderer, w io.Writer) error {
	t := template.New("").Funcs(template.FuncMap{
		"add":    add,
		"render": lc.renderFn,
	})
	// trimspace is important here, otherwise we could have errant Text nodes
	// in the browser screwing up the patch index
	if _, err := t.Parse(strings.TrimSpace(renderer.Render())); err != nil {
		return err
	}
	return t.Execute(w, renderer)
}
func (lc *liveComponent) renderFn(i interface{}) (interface{}, error) {
	renderer, ok := i.(Renderer)
	if !ok {
		return nil, errors.New("cannot render a type that does not implement cave.Renderer")
	}
	var buf strings.Builder

	// check if we already have a renderer with this memory address
	subcomponentAddr := reflect.ValueOf(renderer).Pointer()
	var cr componentRenderer
	cr, ok = lc.subcomponents[subcomponentAddr]
	if !ok {
		// if we don't, assign them an id and store them
		cr = componentRenderer{
			Renderer: renderer,
			id:       lc.counter,
		}
		lc.subcomponents[subcomponentAddr] = cr
		lc.idMap[lc.counter] = subcomponentAddr
		lc.counter++
		// TODO: consider finalizer?
	}
	buf.WriteString(fmt.Sprintf("<div cave-subcomponent=\"%d\">", cr.id))
	if err := lc.render(renderer, &buf); err != nil {
		return nil, err
	}
	buf.WriteString("</div>")
	return buf.String(), nil
}

func renderOnce(r Renderer, w io.Writer) (*liveComponent, error) {
	lc := liveComponent{
		renderer:      r,
		subcomponents: map[uintptr]componentRenderer{},
		idMap:         map[int]uintptr{},
	}
	return &lc, lc.render(r, w)
}

func newLiveComponent(r Renderer) (*liveComponent, error) {
	var buf bytes.Buffer
	lc, err := renderOnce(r, &buf)
	if err != nil {
		return nil, err
	}
	node, err := parseHTMLToNode(buf.String())
	if err != nil {
		return nil, err
	}
	lc.tree = node
	// printNode(lc.tree, 0)
	return lc, nil
}

func (lc *liveComponent) diff() ([]Patch, error) {
	old := lc.tree
	var buf bytes.Buffer
	var err error
	if err = lc.render(lc.renderer, &buf); err != nil {
		return nil, err
	}
	lc.tree, err = parseHTMLToNode(buf.String())
	if err != nil {
		return nil, err
	}
	// printNode(ls.tree, 0)
	return Diff(old, lc.tree)
}

func (wss *websocketSession) sendError(err error) {
	if err == nil {
		return
	}
	msg := messages.ServerMessage{Type: messages.ServerTypeError, Data: []string{err.Error()}}
	b, err := msg.Serialize()
	if err != nil {
		// TODO
		panic(err)
	}

	// TODO: confirm this error is likely just a broken conn and nothing more we should be worried about
	_ = wss.conn.WriteMessage(websocket.TextMessage, b)
}

func (wss *websocketSession) handleMessages() {
	_ = wss.conn.SetReadDeadline(time.Now().Add(pongWait))
	wss.conn.SetPongHandler(func(string) error { _ = wss.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	go func() {
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
			}
		}
	}()
	ticker := time.NewTicker(time.Second * 30)
	for range ticker.C {
		_ = wss.conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := wss.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			return
		}
	}
}

func (wss *websocketSession) handleClientMessage(msg messages.ClientMessage) error {
	spew.Dump(msg)
	switch msg.Type {
	case messages.ClientTypeInit:
		var componentIDs []string
		if err := json.Unmarshal(msg.Data, &componentIDs); err != nil {
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
			var err error
			wss.components[componentID], err = newLiveComponent(renderer)
			if err != nil {
				return err
			}
			onmounter, ok := renderer.(OnMounter)
			if ok {
				onmounter.OnMount(wss.req)
				return wss.reRenderComponent(componentID)
			}
		}
	case messages.ClientTypeSubmit:
		var data []map[string]string
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			return err
		}
		ls, ok := wss.components[msg.ComponentID]
		if !ok {
			return nil
		}
		fmt.Println(data)
		if msg.SubcomponentID == nil {
			ls.renderer.(OnSubmiter).OnSubmit(msg.Name, data[0])
		} else {
			ptr, ok := ls.idMap[*msg.SubcomponentID]
			if !ok {
				log.Fatalf("subcomponent id not found %q", *msg.SubcomponentID)
			}
			ls.subcomponents[ptr].Renderer.(OnSubmiter).OnSubmit(msg.Name, data[0])
		}
		return wss.reRenderComponent(msg.ComponentID)
	case messages.ClientTypeClick:
		ls, ok := wss.components[msg.ComponentID]
		if !ok {
			return nil
		}
		if msg.SubcomponentID == nil {
			ls.renderer.(OnClicker).OnClick(msg.Name)
		} else {
			ptr, ok := ls.idMap[*msg.SubcomponentID]
			if !ok {
				log.Fatalf("subcomponent id not found %q", *msg.SubcomponentID)
			}
			ls.subcomponents[ptr].Renderer.(OnClicker).OnClick(msg.Name)
		}
		return wss.reRenderComponent(msg.ComponentID)
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
		// nothing to do
		return nil
	}
	msg := messages.ServerMessage{
		ComponentID: componentID,
		Type:        messages.ServerTypePatch,
		Data:        diff,
	}
	b, err := msg.Serialize()
	if err != nil {
		return err
	}
	return wss.conn.WriteMessage(websocket.TextMessage, b)
}
