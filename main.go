package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type hnPost struct {
	Title string
	Score int
	URL   string
	ID    string
}

func main() {
	// Get the HN frontpage
	resp, err := http.Get("https://news.ycombinator.com")
	if err != nil {
		fmt.Println("Error getting HN:", err)
		os.Exit(1)
	}

	// Parse the HN frontpage
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("There was an error parsing HN html response:", err)
		os.Exit(1)
	}

	var hnPosts []*hnPost
	hnEntry := &hnPost{}
	doc.Find("table.itemlist tbody tr").Each(func(i int, s *goquery.Selection) {
		if i%3 == 0 {
			ID, _ := s.Attr("id")
			link := s.Find("td.title:last-child a.storylink")
			URL, _ := link.Attr("href")
			hnEntry.URL = string(URL)
			hnEntry.Title = link.Text()
			hnEntry.ID = ID
		} else if i%3 == 1 {
			stats := s.Find("td.subtext")
			score := strings.Split(stats.Find("span.score").Text(), " ")[0]
			hnEntry.Score, _ = strToInt(score)
			if hnEntry.ID != "" {
				hnPosts = append(hnPosts, hnEntry)
			}
			hnEntry = &hnPost{}
		}
	})

	for _, post := range hnPosts {
		fmt.Printf("%+v\n", post)
	}
}

func strToInt(str string) (int, error) {
	nonFractionalPart := strings.Split(str, ".")
	return strconv.Atoi(nonFractionalPart[0])
}
