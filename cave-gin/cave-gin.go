package cavegin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/maxmcd/cave"
)

func MakeRender(handlerContext *HandlerContext, renderer cave.Renderer) render.Render {
	return &Renderer{
		renderer:       renderer,
		handlerContext: handlerContext,
	}
}

type Renderer struct {
	renderer       cave.Renderer
	handlerContext *HandlerContext
}

func (r *Renderer) Render(w http.ResponseWriter) error {
	return r.handlerContext.rc.Render(r.renderer, w)
}

func (r *Renderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/html")
}

type HandlerContext struct {
	rc cave.RenderContext
}

func (hc *HandlerContext) Handler(rendererFunc func() cave.Renderer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Render(200, MakeRender(hc, rendererFunc()))
	}
}
func (hc *HandlerContext) SetTemplateFile(filePath string) error {
	return hc.rc.SetTemplateFile(filePath)
}

func (hc *HandlerContext) SetWebsocketRoute(relativePath string, router gin.IRoutes) {
	hc.rc.SetWebsocketPath(relativePath)

}

func New() *HandlerContext {
	return &HandlerContext{
		rc: cave.RenderContext{},
	}
}
