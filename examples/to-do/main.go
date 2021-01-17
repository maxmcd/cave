package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maxmcd/cave"
	cavegin "github.com/maxmcd/cave/cave-gin"
)

func main() {
	r := gin.Default()
	handlerContext := cavegin.New()
	if err := handlerContext.SetTemplateFile("layout.html"); err != nil {
		log.Fatal(err)
	}
	// handlerContext.
	r.GET("/", handlerContext.Handler(NewToDoApp))
	log.Fatal(r.Run())
}

func NewToDoApp() cave.Renderer {
	tda := &ToDoApp{Items: []string{"breathe"}}
	tld := &ToDoList{}
	tda.ToDoList = tld
	tld.ToDoApp = tda
	return tda
}

type ToDoApp struct {
	Items    []string
	ToDoList *ToDoList
}

var (
	_ cave.OnSubmiter = new(ToDoApp)
	_ cave.Renderer   = new(ToDoApp)
)

func (tda *ToDoApp) OnSubmit(name string, form cave.Form) {

}
func (tda *ToDoApp) Render() string {
	return `
	<div>
	<h3>TODO</h3>
	{{ render .ToDoList }}
	<form cave-submit=todo>
	  <label for="new-todo">
		What needs to be done?
	  </label>
	  <input
		id="new-todo"
	  />
	  <button>
		Add {{ len .Items | add 1 }}
	  </button>
	</form>
  </div>
`
}

type ToDoList struct {
	ToDoApp *ToDoApp
}

func (tdl *ToDoList) Render() string {
	return `
	<ul>
	 {{range .ToDoApp.Items }}
	 	<li>{{.}}</li>
	 {{end}}
	</ul>
`
}

var _ cave.Renderer = new(ToDoList)
