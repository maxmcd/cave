package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maxmcd/cave"
)

func main() {
	r := gin.Default()
	cavern := cave.New()
	if err := cavern.AddTemplateFile("main", "layout.html"); err != nil {
		log.Fatal(err)
	}
	cavern.AddComponent("main", NewToDoApp)
	r.Use(func(c *gin.Context) {
		_, ok := c.Request.URL.Query()["cavews"]
		if ok {
			cavern.ServeWS(c.Writer, c.Request)
			c.Abort()
		}
	})
	r.GET("/", func(c *gin.Context) {
		c.Writer.Header().Add("Content-Type", "text/html")
		if err := cavern.Render(c.Writer); err != nil {
			panic(err)
		}
	})
	r.GET("/bundle.js", func(c *gin.Context) {
		cavern.ServeJS(c.Writer, c.Request)
	})
	r.GET("/bundle.js.map", func(c *gin.Context) {
		cavern.ServeJS(c.Writer, c.Request)
	})
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

func (tda *ToDoApp) OnSubmit(name string, form map[string]string) {
	if name == "todo" {
		tda.Items = append(tda.Items, form["new"])
	}
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
	  <input type="text" name="new" />
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

var _ cave.Renderer = new(ToDoList)

func (tdl *ToDoList) Render() string {
	return `
	<ul>
	 {{range .ToDoApp.Items }}
	 	<li>{{.}}</li>
	 {{end}}
	</ul>
`
}
