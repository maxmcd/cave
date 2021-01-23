package cave

import (
	"bytes"
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
type OnClicker interface {
	OnClick(name string)
}
type Renderer interface {
	Render() string
}

type OnMounter interface {
	OnMount(Session)
}

type Session interface {
	Render()
}

type Cave struct {
	template *template.Template
	registry map[string]RendererFunc
}

func New() *Cave {
	rand.Seed(time.Now().UnixNano())
	return &Cave{}
}

func (cave *Cave) component(name string) (interface{}, error) {
	cmp, ok := cave.registry[name]
	if !ok {
		return nil, fmt.Errorf("component with name %q doesn't exist in the registry", name)
	}
	var buf bytes.Buffer

	// this is a bit ad-hoc. liveview has lots of metadata here. think on this
	fmt.Fprintf(&buf, "<div cave-component=\"%s-%d\">", name, rand.Uint64())
	if _, err := renderOnce(cmp(), &buf); err != nil {
		return nil, err
	}
	fmt.Fprintf(&buf, "</div>")
	return buf.String(), nil
}

func (cave *Cave) AddTemplateFile(name, filePath string) error {
	if cave.template == nil {
		cave.template = template.New(name)
		cave.template.Funcs(template.FuncMap{
			"component": cave.component,
		})
	} else {
		cave.template.New(name)
	}
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if _, err := cave.template.Parse(string(contents)); err != nil {
		return err
	}
	return nil
}

type RendererFunc func() Renderer

func (cave *Cave) AddComponent(name string, rf RendererFunc) {
	if cave.registry == nil {
		cave.registry = map[string]RendererFunc{}
	}
	cave.registry[name] = rf
}

func (cave *Cave) Render(layout string, w io.Writer) error {
	if err := cave.template.ExecuteTemplate(w, layout, nil); err != nil {
		return err
	}
	return nil
}

func (cave *Cave) ServeJS(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Encoding", "gzip")
	w.Header().Add("Content-Type", "application/javascript")
	if strings.HasSuffix(req.URL.Path, ".map") {
		fmt.Fprint(w, bundlejsmap)
	} else {
		fmt.Fprint(w, bundlejs)
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
	_ = wss.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
