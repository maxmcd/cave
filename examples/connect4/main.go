package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maxmcd/cave"
)

func main() {
	r := gin.Default()
	cavern := cave.New()
	if err := cavern.AddTemplateFile("main", "layout.html"); err != nil {
		log.Fatal(err)
	}
	cavern.AddComponent("main", NewConnect4)

	r.Use(func(c *gin.Context) {
		_, ok := c.Request.URL.Query()["cavews"]
		if ok {
			cavern.ServeWS(c.Writer, c.Request)
			c.Abort()
		}
	})
	r.GET("/", func(c *gin.Context) {
		c.Writer.Header().Add("Content-Type", "text/html")
		if err := cavern.Render("main", c.Writer); err != nil {
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

func NewConnect4() cave.Renderer {
	return &Connect4{}
}

type Connect4 struct {
	Username string
	Board    *Board
}

var (
	_ cave.OnSubmiter = new(Connect4)
	_ cave.Renderer   = new(Connect4)
)

func (tda *Connect4) OnSubmit(name string, form map[string]string) {
	if name == "username" {
		tda.Username = form["username"]
		tda.Board = &Board{}
	}
}
func (tda *Connect4) Render() string {
	return `
	<div>
	<h3>Connect4</h3>
	{{ if .Username -}}
	<p><b>username:</b> {{.Username}}</p>
	{{ render .Board }}
	{{- end}}
	{{ if eq .Username "" }}
	<form cave-submit=username>
	  <label for="new-todo">
		Enter your username:
	  </label>
	  <input type="text" name="username" />
	  <input type="submit" />
	</form>
	{{ end }}
  </div>
`
}

type CircleType uint8

const (
	CircleTypeNone CircleType = iota
	CircleTypeRed
	CircleTypeBlack
)

type Board struct {
	board [7][6]CircleType
}

func (board *Board) OnSubmit(name string, form map[string]string) {
	fmt.Println("BOARD ON SUBMIT", name, form)
}

func (board *Board) Render() string {
	var sb strings.Builder
	sb.WriteString(`<div class="board">`)
	for i := 0; i < 7; i++ {
		sb.WriteString(fmt.Sprintf(`<div class="column" cave-click="%d">`, i))
		for j := 0; j < 6; j++ {
			switch board.board[i][j] {
			case CircleTypeRed:
				sb.WriteString(`<div class="red circle"></div>`)
			case CircleTypeBlack:
				sb.WriteString(`<div class="black circle"></div>`)
			default:
				sb.WriteString(`<div class="circle"></div>`)
			}
		}
		sb.WriteString("</div>")
	}
	sb.WriteString("</div>")
	sb.WriteString(`<form cave-submit=board><input type=submit></form>`)
	return sb.String()
}
