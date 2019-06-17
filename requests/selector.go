package requests

import (
	"bytes"
	"gopkg.in/xmlpath.v2"
	"io"
)

type Selector struct {
	body    []byte
	reader  io.Reader
	root    *xmlpath.Node
}

func Node2Selector(node *xmlpath.Node) *Selector {
	return &Selector{
		body: node.Bytes(),
		root: node,
	}
}

func (selector *Selector) Gets(xpath interface{}) ([]*xmlpath.Node, error) {
	if selector.reader == nil {
		selector.reader = bytes.NewReader(selector.body)
	}

	var err error
	root := selector.root
	if root == nil {
		root, err = xmlpath.ParseHTML(selector.reader)
		if err != nil {
			panic(err)
		}
		selector.root = root
	}

	var xpaths []string
	switch xpath.(type) {
	case string:
		xpaths = []string{xpath.(string)}
	case []string:
		xpaths = xpath.([]string)
	}

	for _, xpath := range xpaths {
		result := []*xmlpath.Node{}
		_path := xmlpath.MustCompile(xpath)
		iter := _path.Iter(root)
		for iter.Next() {
			result = append(result, iter.Node())
		}

		if len(result) > 0 {
			return result, nil
		}

	}

	return []*xmlpath.Node{}, NodeNotFound(xpaths)
}

func (selector Selector) Get(xpath interface{}) (*xmlpath.Node, error) {
	node, err := selector.Gets(xpath)
	if err != nil {
		return nil, err
	} else {
		return node[0], err
	}
}

func (selector *Selector) GetNode(xpaths interface{}) (*Selector, error) {
	node, err := selector.Get(xpaths)
	if err != nil {
		return nil, err
	} else {
		return Node2Selector(node), nil
	}
}

func (selector *Selector) GetNodes(xpaths interface{}) ([]*Selector, error) {
	selectors := []*Selector{}
	nodes, err := selector.Gets(xpaths)

	if err != nil {
		return selectors, err
	} else {
		for _, item := range nodes {
			selectors = append(selectors, Node2Selector(item))
		}

		return selectors, nil
	}
}

func (selector *Selector)MustGet(xpath string) string {
	result,err := selector.Get(xpath)
	if err != nil {
		return ""
	}
	return result.String()
}