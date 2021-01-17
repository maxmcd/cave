package cave

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/davecgh/go-spew/spew"
)

type OnSubmiter interface {
	OnSubmit(name string, form Form)
}

type Renderer interface {
	Render() string
}

type Form struct{}

func Render(renderer Renderer, w io.Writer) error {
	// TODO: inject javascript!
	t := template.New("").Funcs(template.FuncMap{
		"add":    add,
		"render": render,
	})
	_, err := t.Parse(renderer.Render())
	if err != nil {
		return err
	}
	return t.Execute(w, renderer)
}

func render(i interface{}) (interface{}, error) {
	renderer, ok := i.(Renderer)
	if !ok {
		return nil, errors.New("cannot render a type that does not implement cave.Renderer")
	}
	var buf bytes.Buffer
	if err := Render(renderer, &buf); err != nil {
		return nil, err
	}
	return buf.String(), nil
}

type RenderContext struct {
	template      *template.Template
	websocketPath string
}

func (rc *RenderContext) SetWebsocketPath(relativePath string) {
	rc.websocketPath = relativePath
}

func (rc *RenderContext) SetTemplateFile(filePath string) error {
	if rc.template != nil {
		return errors.New("can't add a template file twice")
	}
	rc.template = template.New("")

	// provide an empty yield function to satisfy the template.Parse
	rc.template.Funcs(template.FuncMap{
		"yield": func() (interface{}, error) { return nil, nil },
	})
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if _, err := rc.template.Parse(string(contents)); err != nil {
		return err
	}
	return nil
}

func (rc *RenderContext) Render(renderer Renderer, w io.Writer) error {
	yield := func() (interface{}, error) {
		// TODO, cannot call yield more than once
		var buf bytes.Buffer
		if err := Render(renderer, &buf); err != nil {
			return nil, err
		}
		return buf.String(), nil
	}
	rc.template.Funcs(template.FuncMap{
		"yield": yield,
	})
	if err := rc.template.Execute(w, nil); err != nil {
		spew.Dump(err)
		return err
	}
	return nil
}
