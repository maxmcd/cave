package dom

import "strings"

type Node interface {
	Print() string
}

type node struct {
	tag       string
	singleton bool
	children  []interface{}
}

func (n node) Print() string {
	var sb strings.Builder
	sb.WriteString("<")
	sb.WriteString(n.tag)
	if n.singleton {
		sb.WriteString(" />")
		// TODO: error on invalid children
		return sb.String()
	}
	sb.WriteString(">")
	for _, child := range n.children {
		switch v := child.(type) {
		case Node:
			sb.WriteString(v.Print())
		case string:
			sb.WriteString(v)
		}
	}
	sb.WriteString("</")
	sb.WriteString(n.tag)
	sb.WriteString(">")
	return sb.String()
}

func Div(v ...interface{}) Node {
	return node{tag: "div", children: v}
}
func Input(v ...interface{}) Node {
	return node{tag: "input", singleton: true, children: v}
}
func Button(v ...interface{}) Node {
	return node{tag: "button", children: v}
}
func OnClick(v ...interface{}) (i interface{}) {
	return
}

func OnChange(cb func(string)) interface{} {
	cb("hi")
	return nil
}
