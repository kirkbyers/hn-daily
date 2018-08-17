package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kirkbyers/hn-daily/collection"
	"github.com/kirkbyers/hn-daily/db"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Something went wrong finding homeDir:", err)
		os.Exit(1)
	}
	dbPath := filepath.Join(home, "hn-daily.db")
	db.Init(dbPath)

	t := time.NewTicker(15 * time.Minute)
	defer t.Stop()

	collectPosts()
	for {
		select {
		case <-t.C:
			collectPosts()
		}
	}
}

func collectPosts() {
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
