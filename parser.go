package main

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getTextValue(n *html.Node) string {
	for v := range n.Descendants() {
		if v.Type == html.TextNode {
			return v.Data
		}
	}

	return ""
}

func getAttribute(n *html.Node, key string) string {
	for _, v := range n.Attr {
		if v.Key == key {
			return v.Val
		}
	}

	return ""
}

func isTag(n *html.Node, a atom.Atom) bool {
	if n.Type == html.ElementNode && n.DataAtom == a {
		return true
	}

	return false
}

func findByTag(n *html.Node, a atom.Atom) *html.Node {
	if isTag(n, a) {
		return n
	}

	for next := n.FirstChild; next != nil; next = next.NextSibling {
		res := findByTag(next, a)
		if res != nil {
			return res
		}
	}

	return nil
}

func findByClass(n *html.Node, target string) *html.Node {
	if n.Type == html.ElementNode {
		if getAttribute(n, "class") == target {
			return n
		}
	}

	for next := n.FirstChild; next != nil; next = next.NextSibling {
		res := findByClass(next, target)
		if res != nil {
			return res
		}
	}

	return nil
}
