package cave

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"unsafe"

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

type TestComponentParent struct {
	Count1 *TestComponent
	Count2 *TestComponent
	Count3 Renderer // will this behave differently?

	order bool
}

var (
	_ Renderer     = new(TestComponentParent)
	_ RendererFunc = NewTestComponentParent
)

func NewTestComponentParent() Renderer {
	return &TestComponentParent{
		Count1: &TestComponent{},
		Count2: &TestComponent{},
		Count3: NewTestComponent(),
	}
}
func (tc *TestComponentParent) Render() string {
	if tc.order {
		return `
		{{render .Count1 }}
		{{render .Count2 }}
		{{render .Count3 }}
		`
	}
	return `
	{{render .Count3 }}
	{{render .Count2 }}
	{{render .Count1 }}
	`
}

func TestCaveBasic(t *testing.T) {
	cavern := New()
	if err := cavern.AddTemplateFile("main", "./examples/to-do/layout.html"); err != nil {
		t.Fatal(err)
	}
	cavern.AddComponent("main", NewTestComponent)

	var buf bytes.Buffer

	_ = cavern.Render("main", &buf)
	fmt.Println(buf.String())
}

func TestCaveDiff(t *testing.T) {
	tc := NewTestComponent()
	ls, err := newLiveComponent(tc, nil)
	if err != nil {
		t.Fatal(err)
	}
	tc.(OnSubmiter).OnSubmit("name", nil)
	spew.Dump(ls.diff())
}

func TestPointerTricks(t *testing.T) {
	tc := NewTestComponent()
	tc1 := NewTestComponent()
	tc2 := NewTestComponent()

	fmt.Println(reflect.ValueOf(tc).Pointer())

	store := map[uintptr]struct{}{}

	store[reflect.ValueOf(tc).Pointer()] = struct{}{}

	_, ok := store[uintptr(unsafe.Pointer(&tc))]
	fmt.Println(tc1, tc2, ok)
}
