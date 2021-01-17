package cave

import (
	"bytes"
	"encoding/json"
	"reflect"

	"golang.org/x/net/html"
)

type PatchType uint8

const (
	PatchTypeInsert PatchType = iota
	PatchTypeRemove
	PatchTypeAttributes
	PatchTypeText
	PatchTypeElement
)

type Patch struct {
	Type       PatchType     `json:"t"`
	Data       string        `json:"d,omitempty"`
	Attributes AttributeList `json:"a,omitempty"`
	Index      int           `json:"i"`
	node       *html.Node
}

type AttributeList []Attribute

func (a AttributeList) MarshalJSON() ([]byte, error) {
	out := [][]string{}
	for _, at := range a {
		out = append(out, []string{at.Key, at.Val})
	}
	return json.Marshal(out)
}

func (p *Patch) UnmarshalJSON(data []byte) error {
	type Tmp Patch
	tmp := &struct {
		Attr [][]string `json:"a"`
		*Tmp
	}{
		Tmp: (*Tmp)(p),
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	for _, tuple := range tmp.Attr {
		p.Attributes = append(p.Attributes, Attribute{Key: tuple[0], Val: tuple[1]})
	}
	return nil
}

type Attribute struct {
	// We don't support namespaces yet
	Key string
	Val string
}

func Apply(a *html.Node, patches []Patch) {
	_, _ = apply(a, patches, 0)
}
func apply(a *html.Node, patches []Patch, index int) ([]Patch, int) {
	if len(patches) == 0 {
		return nil, index
	}
	patch := patches[0]
	if patch.Index == index {
		patches = patches[1:]
		switch patch.Type {
		case PatchTypeRemove:
			a.Parent.RemoveChild(a)
		case PatchTypeAttributes:
			a.Attr = []html.Attribute{}
			for _, attr := range patch.Attributes {
				a.Attr = append(a.Attr, html.Attribute{
					Key:       attr.Key,
					Val:       attr.Val,
					Namespace: a.Namespace,
				})
			}
		case PatchTypeText:
			a.Data = patch.Data
		case PatchTypeElement:
			a.Type = html.RawNode
			a.Data = patch.Data
			a.FirstChild = nil
			a.LastChild = nil
		}
	}

	if len(patches) == 0 {
		return nil, 0
	}
	patch = patches[0]
	if patch.Type == PatchTypeInsert && index+1 == patch.Index {
		a.Parent.AppendChild(&html.Node{Type: html.RawNode, Data: patch.Data})
		patches = patches[1:]
	}

	aChild := a.FirstChild
	for aChild != nil {
		patches, index = apply(aChild, patches, index+1)
		aChild = aChild.NextSibling
	}
	return patches, index
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
		patches[i].Data = buf.String()
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
		return []Patch{{Type: PatchTypeInsert, node: b, Index: index}}, index
	}
	if b == nil {
		return []Patch{{Type: PatchTypeRemove, Index: index}}, index
	}
	if a.Type == html.ElementNode && b.Type == html.ElementNode {
		if a.Data == b.Data {
			// we're all equal, check prop inequality
			if !reflect.DeepEqual(a.Attr, b.Attr) {
				return []Patch{{Type: PatchTypeAttributes, Attributes: attributesToAttributes(b.Attr), Index: index}}, index
			}
		} else {
			// different tag name, replace whole node
			return []Patch{{Type: PatchTypeElement, node: b, Index: index}}, index
		}
	}
	if a.Type == html.TextNode && b.Type == html.TextNode {
		if a.Data != b.Data {
			return []Patch{{Type: PatchTypeText, Data: b.Data, Index: index}}, index
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
