package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	htmlfile := "./ex4.html"
	fh, err := os.Open(htmlfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fh.Close()

	reader := bufio.NewReader(fh)
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	doclink := []Link{}
	var nAnchor *html.Node

	// for every ahref node there are a series of text nodes corresponding to the ahreaf node
	// create a new link whenever an a href is encountered and add it to list of doclinks
	// whenever text is enconutered as a child of anchor node, update the last element of the doclinks and do not create new links
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// start processing the anchor node
			nAnchor = n
			for _, a := range n.Attr {
				if a.Key == "href" {
					newlink := Link{}
					newlink.Href = a.Val
					newlink.Text = ""
					doclink = append(doclink, newlink)
				}
			}
		}

		// get text within current node only if the current
		// is a child of anchor node.
		if nAnchor != nil {
			if n.Type == html.TextNode {
				text := strings.TrimSpace(n.Data)
				if text != "" && len(doclink) != 0 {
					doclink[len(doclink)-1].Text += text
				}
			}
		}

		// get all children of current node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		// all children of current anchor node is processed
		// stop consuming text till new anchor node is encountered
		if n == nAnchor {
			nAnchor = nil
		}
	}

	f(doc)
	fmt.Printf("%v\n", doclink)
}

// TraverseNode: recursively traverse the node and print on screen
func TraverseNode(n *html.Node, depth int) {
	indent := strings.Repeat(" ", depth)

	switch n.Type {
	case html.ElementNode:
		fmt.Printf("%s<%s>\n", indent, n.Data)
	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		fmt.Printf("%s%s\n", indent, text)
	case html.CommentNode:
		fmt.Printf("%s<!-- %s -->\n", indent, n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseNode(c, depth+1)
	}

}
