package cave

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/maxmcd/cave/dom"
)

type Server struct {
	new      func() Renderable
	sessions map[string]Renderable
}

func newSession() (out string, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	return fmt.Sprintf("%x", b), nil
}

func (s Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			sessionKey, err := newSession()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			renderable := s.new()
			s.sessions[sessionKey] = renderable
			node := renderable.Render()
			fmt.Fprint(w, node.Print())
			// write node to body
		} else if req.Method == http.MethodPost {
			// handle state change
		}
	}))
}

func New(new func() Renderable) Server {
	return Server{new: new, sessions: map[string]Renderable{}}
}

type Renderable interface {
	Render() dom.Node
}
