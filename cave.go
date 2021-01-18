package cave

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

type OnSubmiter interface {
	OnSubmit(name string, form map[string]string)
}

type Renderer interface {
	Render() string
}

// MAYBE
type OnMounter interface {
	OnMount(req *http.Request)
}

type Form struct{}

func Render(renderer Renderer, w io.Writer) error {
	t := template.New("").Funcs(template.FuncMap{
		"add":    add,
		"render": render,
	})
	// trimspace is important here, otherwise we could have errant Text nodes
	// in the browser screwing up the patch index
	_, err := t.Parse(strings.TrimSpace(renderer.Render()))
	if err != nil {
		return err
	}
	return t.Execute(w, renderer)
}

func render(i interface{}) (interface{}, error) {
	renderer, ok := i.(Renderer)
	if !ok {
		return nil, errors.New("cannot render a type that does not implement cave.Renderer")
	}
	var buf bytes.Buffer
	if err := Render(renderer, &buf); err != nil {
		return nil, err
	}
	return buf.String(), nil
}

type Cave struct {
	template *template.Template
	registry map[string]RendererFunc
}

func New() *Cave {
	rand.Seed(time.Now().UnixNano())
	return &Cave{}
}

func (rc *Cave) component(name string) (interface{}, error) {
	cmp, ok := rc.registry[name]
	if !ok {
		return nil, fmt.Errorf("component with name %q doesn't exist in the registry", name)
	}
	var buf bytes.Buffer

	// this is a bit ad-hoc. liveview has lots of metadata here. think on this
	fmt.Fprintf(&buf, "<div cave-component=\"%s-%d\">", name, rand.Uint64())
	if err := Render(cmp(), &buf); err != nil {
		return nil, err
	}
	fmt.Fprintf(&buf, "</div>")
	return buf.String(), nil
}

func (rc *Cave) AddTemplateFile(name, filePath string) error {
	if rc.template == nil {
		rc.template = template.New(name)
		rc.template.Funcs(template.FuncMap{
			"component": rc.component,
		})
	} else {
		rc.template.New(name)
	}
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if _, err := rc.template.Parse(string(contents)); err != nil {
		return err
	}
	return nil
}

type RendererFunc func() Renderer

func (rc *Cave) AddComponent(name string, rf RendererFunc) {
	if rc.registry == nil {
		rc.registry = map[string]RendererFunc{}
	}
	rc.registry[name] = rf
}

func (rc *Cave) Render(w io.Writer) error {
	if err := rc.template.Execute(w, nil); err != nil {
		return err
	}
	return nil
}

func (cave *Cave) ServeJS(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Encoding", "gzip")
	w.Header().Add("Content-Type", "application/javascript")
	if strings.HasSuffix(req.URL.Path, ".map") {
		fmt.Fprint(w, bundle_js_map)
	} else {
		fmt.Fprint(w, bundle_js)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (cave *Cave) ServeWS(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	wss := &websocketSession{cave: cave, req: req, conn: conn}
	wss.handleMessages()
}
