package collection

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kirkbyers/hn-daily/db"
)

// Collect scrapes the front page of HN
func Collect() ([]db.HnPost, error) {
	doc, err := getHnFrontPage()
	if err != nil {
		return nil, err
	}

	return parseHnFPDoc(doc), nil
}

func getHnFrontPage() (*goquery.Document, error) {
	// Get the HN frontpage
	resp, err := http.Get("https://news.ycombinator.com")
	if err != nil {
		fmt.Println("Error getting HN:", err)
		return nil, err
	}

	// Parse the HN frontpage from http.Response to goquery doc
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("There was an error parsing HN html response:", err)
		return nil, err
	}
	return doc, nil
}

// Parse goquery doc to usable data structure
func parseHnFPDoc(doc *goquery.Document) (hnPosts []db.HnPost) {
	hnEntry := db.HnPost{}

	// Find the table rows that make up the front page
	doc.Find("table.itemlist tbody tr").Each(func(i int, s *goquery.Selection) {
		/***************************
		* A post on the frontpage is split into 3 consecutive rows in the table
		* 0: postId is in id attr, URL and postTitle are in an anchor
		* 1: score is in a span here
		* 2: Spacer row
		* Posts' triplet pairs are consecutive for the entirety of the table
		***************************/
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
			hnEntry = db.HnPost{}
		}
	})
	return hnPosts
}

// Helper func to parse strings to ints
func strToInt(str string) (int, error) {
	nonFractionalPart := strings.Split(str, ".")
	return strconv.Atoi(nonFractionalPart[0])
}
