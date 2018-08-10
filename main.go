package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	// Reused var
	var err error

	// Get the HN frontpage
	var resp *http.Response
	resp, err = http.Get("https://news.ycombinator.com")
	if err != nil {
		fmt.Println("Error getting HN:", err)
		os.Exit(1)
	}

	// Parse the HN frontpage
	var doc *html.Node
	doc, err = html.Parse(resp.Body)
	if err != nil {
		fmt.Println("There was an error parsing HN html response:", err)
		os.Exit(1)
	}

	// Search for table of posts
	var t *html.Node
	t, err = searchNodeForClass(doc.FirstChild, "itemlist")
	if err != nil {
		fmt.Println("There was an issue searching the html")
		os.Exit(1)
	}

	posts, _ := parseHNTable(t)
	fmt.Printf("%+v\n", len(posts))
}

type hnPost struct {
	Title    string
	Score    string
	Comments string
	URL      string
	ID       string
}

func parseHNTable(t *html.Node) ([]*hnPost, error) {
	var result []*hnPost
	var tBody *html.Node
	for tBody = t.FirstChild; tBody.Type != html.ElementNode; tBody = tBody.NextSibling {
	}
tableLoop:
	for c := tBody.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue tableLoop
		}
		hnEntry := &hnPost{}
		// Check attrs
		for _, a := range c.Attr {
			if a.Key == "class" {
				if a.Val == "spacer" {
					continue tableLoop
				} else if a.Val == "athing" {
					ID, title, URL, err := processAThing(c)
					if err != nil {
						fmt.Println(err)
					}
					hnEntry.ID = ID
					hnEntry.Title = title
					hnEntry.URL = URL
				}
			}
		}
		if len(c.Attr) == 0 {
			fmt.Printf("%+v\n", c.LastChild)
		}
		result = append(result, hnEntry)
	}
	return result, nil
}

// func processSubRow(s *html.Node)

func processAThing(a *html.Node) (ID, title, URL string, err error) {
	for _, attr := range a.Attr {
		if attr.Key == "id" {
			ID = attr.Val
		}
	}

	for c := a.FirstChild; c != nil; c = c.NextSibling {
		if c.NextSibling != nil || c.Type != html.ElementNode {
			continue
		}
		link := c.FirstChild
		for _, a := range link.Attr {
			if a.Key == "href" {
				URL = a.Val
			}
		}
		title = link.FirstChild.Data
	}
	if title == "" || URL == "" {
		return ID, title, URL, errors.New("Unable to extract title or URL from \"athing\" node")
	}
	return ID, title, URL, nil
}

func searchNodeForClass(n *html.Node, s string) (*html.Node, error) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		for _, a := range c.Attr {
			if a.Key == "class" && a.Val == s {
				return c, nil
			}
		}
		res, err := searchNodeForClass(c, s)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, nil
}
