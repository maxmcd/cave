package cavegin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/maxmcd/cave"
)

func MakeRender(cavern *Cave, renderer cave.Renderer) render.Render {
	return &Renderer{
		renderer: renderer,
		cavern:   cavern,
	}
}

type Renderer struct {
	renderer cave.Renderer
	cavern   *Cave
}

func (r *Renderer) Render(w http.ResponseWriter) error {
	return r.cavern.Render(r.renderer, w)
}

func (r *Renderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/html")
}

type Cave struct {
	cave.Cave
}

func (hc *Cave) Handler(rendererFunc func() cave.Renderer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Render(200, MakeRender(hc, rendererFunc()))
	}
}

func New() *Cave {
	return &Cave{
		Cave: cave.Cave{},
	}
}
