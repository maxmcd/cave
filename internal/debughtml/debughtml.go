package debughtml

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func PrettyNode(a *html.Node) string {
	var buf bytes.Buffer
	switch a.Type {
	case html.TextNode:
		return fmt.Sprintf("text: %q", a.Data)
	case html.ElementNode:
		return fmt.Sprintf("element: %s %s", a.Data, a.Attr)
	default:
		_ = html.Render(&buf, a)
		return "unknown: " + buf.String()
	}
}

func PrintNode(a *html.Node, level int) {
	fmt.Println(strings.Repeat("\t", level), PrettyNode(a))
	for _, child := range NodeChildren(a) {
		PrintNode(child, level+1)
	}
}

func NodeChildren(a *html.Node) (children []*html.Node) {
	child := a.FirstChild
	for child != nil {
		children = append(children, child)
		child = child.NextSibling
	}
	return
}
