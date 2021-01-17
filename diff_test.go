package cave

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/html"
)

// func Test_walk(t *testing.T) {
// 	spew.Dump(walk(mustParse("<div><div></div></div>"), mustParse("<div></div>"), 0))
// }

func mustParse(input string) *html.Node {
	node, err := html.Parse(bytes.NewBuffer([]byte(input)))
	if err != nil {
		panic(err)
	}
	// it goes <html><head></head><body>, so just return the body
	return node.FirstChild.FirstChild.NextSibling
}

func Test_Diff(t *testing.T) {
	tests := []struct {
		name        string
		a           string
		b           string
		wantPatches []Patch
	}{
		{
			name:        "add div",
			a:           "<div></div>",
			b:           "<div></div><div></div>",
			wantPatches: []Patch{{pt: PatchTypeInsert, data: "<div></div>", index: 2}},
		},
		{
			name:        "add div",
			a:           "<div></div><div></div>",
			b:           "<div></div>",
			wantPatches: []Patch{{pt: PatchTypeRemove, index: 2}},
		},
		{
			name: "further nested, no patches",
			a:    `<form><input type="text"/><input type="submit"/></form>`,
			b:    `<form><input type="text"/><div>Extra</div><input type="submit"/></form>`,
		},
		{
			name: "further nested, no patches",
			a:    `<form><input type="text"/><input type="submit"/></form>`,
			b:    `<form><input type="text"/><input type="submit"/><div>Extra</div></form>`,
		},
		{
			name: "further nested, no patches",
			a:    `<form><input type="text"/><input type="submit"/></form>`,
			b:    `<form><div>Extra</div><input type="text"/><input type="submit"/></form>`,
		},
		{
			name: "replace tags",
			a:    "<div foo=bar one=two>Hello</div>",
			b:    `<div foo="baz" three="four">Hello</div>`,
			wantPatches: []Patch{{pt: PatchTypeAttributes, attributes: []Attribute{
				{Key: "foo", Val: "baz"},
				{Key: "three", Val: "four"},
			}, index: 1}},
		},
		{
			name:        "replace words",
			a:           "<div>Hello</div>",
			b:           "<div>World</div>",
			wantPatches: []Patch{{pt: PatchTypeText, data: "World", index: 2}},
		},
		{
			name: "swap tag text",
			a:    "<div>Hello</div><div>World</div>",
			b:    "<div>World</div><div>Hello</div>",
			wantPatches: []Patch{
				{pt: PatchTypeText, data: "World", index: 2},
				{pt: PatchTypeText, data: "Hello", index: 4},
			},
		},
		{
			name: "swap element",
			a:    "<div>Hello</div>",
			b:    `<span foo="bar">Hello</span>`,
			wantPatches: []Patch{
				{pt: PatchTypeElement, data: "<span foo=\"bar\">Hello</span>", index: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := mustParse(tt.a)
			gotPatches, err := Diff(a, mustParse(tt.b))
			if err != nil {
				t.Error(err)
			}
			if tt.wantPatches != nil && !reflect.DeepEqual(gotPatches, tt.wantPatches) {
				t.Errorf("Diff() = %s, want %s",
					spew.Sdump(gotPatches),
					spew.Sdump(tt.wantPatches),
				)
			}
			Apply(a, gotPatches)

			a.Type = html.DocumentNode // hack away the body tag
			a.Data = ""

			var buf bytes.Buffer
			_ = html.Render(&buf, a)
			if buf.String() != tt.b {
				t.Errorf("Diff() = %v, want %v",
					buf.String(),
					tt.b)
			}
		})
	}
}
