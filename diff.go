package cave

import (
	"bytes"
	"reflect"

	"golang.org/x/net/html"
)

type PatchType uint8

const (
	PatchTypeNone PatchType = iota
	PatchTypeInsert
	PatchTypeRemove
	PatchTypeAttributes
	PatchTypeText
	PatchTypeElement
)

type Patch struct {
	pt         PatchType
	data       string
	attributes []Attribute
	node       *html.Node
	index      int
}

type Attribute struct {
	// We don't support namespaces yet
	Key string
	Val string
}

func Apply(a *html.Node, patches []Patch) {
	_, _ = apply(a, patches, 0)
}
func apply(a *html.Node, patches []Patch, index int) (int, []Patch) {
	if len(patches) == 0 {
		return index, nil
	}

	patch := patches[0]
	if patch.index == index {
		patches = patches[1:]
		switch patch.pt {
		case PatchTypeRemove:
			a.Parent.RemoveChild(a)
		case PatchTypeAttributes:
			a.Attr = []html.Attribute{}
			for _, attr := range patch.attributes {
				a.Attr = append(a.Attr, html.Attribute{
					Key:       attr.Key,
					Val:       attr.Val,
					Namespace: a.Namespace,
				})
			}
		case PatchTypeText:
			a.Data = patch.data
		case PatchTypeElement:
			a.Type = html.RawNode
			a.Data = patch.data
			a.FirstChild = nil
			a.LastChild = nil
		}
	}

	if len(patches) == 0 {
		return 0, nil
	}
	patch = patches[0]
	if patch.pt == PatchTypeInsert && index+1 == patch.index {
		a.Parent.AppendChild(&html.Node{Type: html.RawNode, Data: patch.data})
		patches = patches[1:]
	}

	aChild := a.FirstChild
	for aChild != nil {
		index, patches = apply(aChild, patches, index+1)
		aChild = aChild.NextSibling
	}
	return index, patches
}

func Diff(a *html.Node, b *html.Node) ([]Patch, error) {
	patches, _ := walk(a, b, 0)
	var buf bytes.Buffer
	for i, patch := range patches {
		if patch.node == nil {
			continue
		}
		if err := html.Render(&buf, patch.node); err != nil {
			return nil, err
		}
		patches[i].data = buf.String()
		patches[i].node = nil
		buf.Reset()
	}
	return patches, nil
}

func attributesToAttributes(a []html.Attribute) []Attribute {
	out := []Attribute{}
	for _, attr := range a {
		out = append(out, Attribute{Key: attr.Key, Val: attr.Val})
	}
	return out
}

func walk(a *html.Node, b *html.Node, index int) ([]Patch, int) {
	if reflect.DeepEqual(a, b) {
		return nil, index
	}
	if a == nil {
		return []Patch{{pt: PatchTypeInsert, node: b, index: index}}, index
	}
	if b == nil {
		return []Patch{{pt: PatchTypeRemove, index: index}}, index
	}
	if a.Type == html.ElementNode && b.Type == html.ElementNode {
		if a.Data == b.Data {
			// we're all equal, check prop inequality
			if !reflect.DeepEqual(a.Attr, b.Attr) {
				return []Patch{{pt: PatchTypeAttributes, attributes: attributesToAttributes(b.Attr), index: index}}, index
			}
		} else {
			// different tag name, replace whole node
			return []Patch{{pt: PatchTypeElement, node: b, index: index}}, index
		}
	}
	if a.Type == html.TextNode && b.Type == html.TextNode {
		if a.Data != b.Data {
			return []Patch{{pt: PatchTypeText, data: b.Data, index: index}}, index
		}
	}
	patches := []Patch{}
	aChild := a.FirstChild
	bChild := b.FirstChild
	for !(aChild == nil && bChild == nil) {
		var p []Patch
		p, index = walk(aChild, bChild, index+1)
		patches = append(patches, p...)
		if aChild != nil {
			aChild = aChild.NextSibling
		}
		if bChild != nil {
			bChild = bChild.NextSibling
		}
	}
	return patches, index
}
