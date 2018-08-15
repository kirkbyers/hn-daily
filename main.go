package main

import (
	"fmt"
	"os"

	"github.com/kirkbyers/hn-daily/collection"
	"github.com/kirkbyers/hn-daily/db"
)

func main() {
	db.Init()
	posts, err := collection.Collect()
	if err != nil {
		fmt.Println("Something went wrong collecting posts:", err)
		os.Exit(1)
	}
	if len(posts) <= 0 {
		fmt.Println("No posts were recored")
		return
	}
	if err := db.SaveHnScrape(posts); err != nil {
		fmt.Println("Something went wrong saving posts to db:", err)
		os.Exit(1)
	}
}
