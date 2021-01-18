package cave

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type TestComponent struct {
	Count int
}

var (
	_ Renderer     = new(TestComponent)
	_ OnSubmiter   = new(TestComponent)
	_ RendererFunc = NewTestComponent
)

func NewTestComponent() Renderer {
	return &TestComponent{Count: 5}
}

func (tc *TestComponent) Render() string {
	return "<div>{{.Count}}</div>"
}

func (tc *TestComponent) OnSubmit(name string, f map[string]string) {
	tc.Count++
}

func TestCaveBasic(t *testing.T) {
	cavern := New()
	if err := cavern.AddTemplateFile("main", "./examples/to-do/layout.html"); err != nil {
		t.Fatal(err)
	}
	cavern.AddComponent("main", NewTestComponent)

	var buf bytes.Buffer

	cavern.Render(&buf)
	fmt.Println(buf.String())
}

func TestCaveDiff(t *testing.T) {
	tc := NewTestComponent()
	ls := newLiveState(tc)
	tc.(OnSubmiter).OnSubmit("name", nil)
	spew.Dump(ls.diff())
}
