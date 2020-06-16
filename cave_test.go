package cave

import (
	"fmt"
	"html/template"
	"os"
	"testing"
	"text/template/parse"

	"github.com/davecgh/go-spew/spew"
)

func TestRenderWithInt(t *testing.T) {

	tmp, err := template.New("foo").Parse(`
<li class="todo-list-item {{ if .Completed }} completed {{ end }}">
<input class="" type="text" {{ if .Completed }} checked {{ end }}>
<label class="todo-label">{{ .Title }}</label>
<button class="destroy"></button>
<input class="edit" onChange="{{ .OnChange }}" value="{{ .Title }}">
</li>
`)

	spew.Config.DisableMethods = true

	// spew.Dump(tmp.Tree.Root.Nodes[7].(*parse.ActionNode).Pipe.Cmds)
	tmp.Execute(os.Stdout, struct {
		Completed bool
		Title     string
		OnChange  func()
	}{
		Completed: true,
		Title:     "hi",
		OnChange: func() {
			fmt.Println("ok then")
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(tmp.Tree.Root.Nodes[7].(*parse.ActionNode).Pipe.Cmds[0])

}

type Foo func()

func (f Foo) String() string {
	return "hi"
}
