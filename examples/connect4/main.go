package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maxmcd/cave"
)

func main() {

	gm := GameMaster{}
	gm.startGameMaster()

	r := gin.Default()
	cavern := cave.New()
	if err := cavern.AddTemplateFile("main", "layout.html"); err != nil {
		log.Fatal(err)
	}
	cavern.AddComponent("main", gm.NewConnect4)

	r.Use(func(c *gin.Context) {
		if _, ok := c.Request.URL.Query()["cavews"]; ok {
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
	r.GET("/bundle.js", func(c *gin.Context) { cavern.ServeJS(c.Writer, c.Request) })
	r.GET("/bundle.js.map", func(c *gin.Context) { cavern.ServeJS(c.Writer, c.Request) })
	log.Fatal(r.Run())
}

type GameMaster struct {
	playerRequests chan *Connect4
}

func (gm *GameMaster) startGameMaster() {
	gm.playerRequests = make(chan *Connect4)
	go func() {
		for {
			first := <-gm.playerRequests
			second := <-gm.playerRequests
			fmt.Println("got two players", first, second)
			first.Opponent = second
			second.Opponent = first
			side := rand.Intn(1)
			first.Board = &BoardView{parent: first}
			second.Board = &BoardView{parent: second}
			if side == 1 {
				first.Board.Side = CircleTypeBlack
				second.Board.Side = CircleTypeRed
			} else {
				first.Board.Side = CircleTypeRed
				second.Board.Side = CircleTypeBlack
			}
			game := &Game{}
			first.Board.game = game
			second.Board.game = game
			first.session.Render()
			second.session.Render()
		}
	}()
}
func (gm *GameMaster) NewConnect4() cave.Renderer {
	return &Connect4{
		playerRequests: gm.playerRequests,
	}
}

type Connect4 struct {
	Username       string
	Board          *BoardView
	session        cave.Session
	Opponent       *Connect4
	playerRequests chan *Connect4
}

var (
	_ cave.OnSubmiter = new(Connect4)
	_ cave.Renderer   = new(Connect4)
)

func (tda *Connect4) OnMount(session cave.Session) {
	tda.session = session
}

func (tda *Connect4) OnSubmit(name string, form map[string]string) {
	if name == "username" {
		tda.Username = form["username"]
		tda.playerRequests <- tda
	}
}
func (tda *Connect4) Render() string {
	return `
	<div>
	<h3>Connect4</h3>
	{{ if .Username -}}
	<p><b>username:</b> {{.Username}}</p>
		{{ if .Board }}
		Playing against {{ .Opponent.Username }}
		{{ render .Board }}
		{{ else }}
		Waiting for another player to join...
		{{- end}}
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

type BoardView struct {
	game   *Game
	Side   CircleType
	parent *Connect4
}

var (
	_ cave.OnClicker = new(BoardView)
)

func (board *BoardView) OnClick(name string) {
	column, _ := strconv.Atoi(name)
	_ = board.game.play(board.Side, column)
	board.parent.Opponent.session.Render()
}

func (board *BoardView) Render() string {
	var sb strings.Builder
	sb.WriteString(`<div class="board">`)
	for i := 0; i < 7; i++ {
		sb.WriteString(fmt.Sprintf(`<div class="column" cave-click="%d">`, i))
		for j := 5; j >= 0; j-- {
			switch board.game.board[i][j] {
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
	return sb.String()
}

const (
	BoardWidth  = 7
	BoardHeight = 6
)

type Game struct {
	board  [BoardWidth][BoardHeight]CircleType
	winner CircleType
}

func (g *Game) moveNumber() int {
	var count int
	for _, column := range g.board {
		for _, square := range column {
			if square != CircleTypeNone {
				count++
			}
		}
	}
	return count
}
func (g *Game) whosMove() CircleType {
	if g.moveNumber()%2 == 0 {
		return CircleTypeRed
	}
	return CircleTypeBlack
}

func (g *Game) didTheyWin(circle CircleType) bool {
	// Taken from: https://stackoverflow.com/a/38211417/1333724
	for j := 0; j < BoardHeight-3; j++ {
		for i := 0; i < BoardWidth; i++ {
			if g.board[i][j] == circle && g.board[i][j+1] == circle && g.board[i][j+2] == circle && g.board[i][j+3] == circle {
				return true
			}
		}
	}
	for i := 0; i < BoardWidth-3; i++ {
		for j := 0; j < BoardHeight; j++ {
			if g.board[i][j] == circle && g.board[i+1][j] == circle && g.board[i+2][j] == circle && g.board[i+3][j] == circle {
				return true
			}
		}
	}
	for i := 3; i < BoardWidth; i++ {
		for j := 0; j < BoardHeight-3; j++ {
			if g.board[i][j] == circle && g.board[i-1][j+1] == circle && g.board[i-2][j+2] == circle && g.board[i-3][j+3] == circle {
				return true
			}
		}
	}
	for i := 3; i < BoardWidth; i++ {
		for j := 3; j < BoardHeight; j++ {
			if g.board[i][j] == circle && g.board[i-1][j-1] == circle && g.board[i-2][j-2] == circle && g.board[i-3][j-3] == circle {
				return true
			}
		}
	}
	return false
}

func (g *Game) play(circle CircleType, column int) (winner CircleType) {
	if g.winner != CircleTypeNone {
		return
	}
	if circle != g.whosMove() {
		return
	}
	for i, square := range g.board[column] {
		if square == CircleTypeNone {
			g.board[column][i] = circle
			if g.didTheyWin(circle) {
				g.winner = circle
				return circle
			}
			break
		}
	}
	return CircleTypeNone
}
