package main

import (
	"fmt"

	"github.com/maxmcd/cave"
)

type App struct {
	Counter    *int
	SetCounter cave.IntHook
}

func Init(hooks cave.Hooks) cave.Renderer {
	app := &App{}
	app.Counter, app.SetCounter = hooks.NewIntHook()
	return app
}

func (app *App) Render(ctx cave.Context, ui cave.UI) (err error) {
	ui.Div(func() {
		ui.Text(fmt.Sprint("Hello World", *app.Counter))
	}, ctx.OnClick(app.ClickHello))
	return nil
}

func (app *App) ClickHello(event cave.Event) {
	app.SetCounter(*app.Counter + 1)
}

func main() {

}
