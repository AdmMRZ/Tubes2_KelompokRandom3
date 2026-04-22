package parser

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Node struct {
	Tag        string            `json:"Tag"`
	Attributes map[string]string `json:"Attributes"`
	Children   []*Node           `json:"Children"`
	Parent     *Node             `json:"-"`
	IsText     bool              `json:"-"`
}

func parseReader(r io.Reader) (*Node, error) {
	n, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var walk func(hn *html.Node, parent *Node) *Node
	walk = func(hn *html.Node, parent *Node) *Node {
		if hn.Type == html.ElementNode && (hn.Data == "script" || hn.Data == "style" || hn.Data == "noscript") {
			return nil
		}
		if hn.Type == html.TextNode && strings.TrimSpace(hn.Data) == "" {
			return nil
		}

		res := &Node{
			Tag:        strings.TrimSpace(hn.Data),
			Attributes: make(map[string]string),
			Parent:     parent,
			IsText:     hn.Type == html.TextNode,
		}

		for _, a := range hn.Attr {
			res.Attributes[a.Key] = a.Val
		}

		for c := hn.FirstChild; c != nil; c = c.NextSibling {
			if child := walk(c, res); child != nil {
				res.Children = append(res.Children, child)
			}
		}

		return res
	}

	return walk(n, nil), nil
}

func ParseHTML(url string) (*Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseReader(resp.Body)
}

func ParseHTMLText(htmlContent string) (*Node, error) {
	return parseReader(strings.NewReader(htmlContent))
}

func MaxDepth(n *Node) int {
	if n == nil {
		return 0
	}
	maxChild := 0
	for _, c := range n.Children {
		d := MaxDepth(c)
		if d > maxChild {
			maxChild = d
		}
	}
	return maxChild + 1
}
