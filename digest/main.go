package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/kirkbyers/hn-daily/db"
	"github.com/kirkbyers/hn-daily/process"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Something went wrong finding homeDir:", err)
		os.Exit(1)
	}
	dbPath := filepath.Join(home, "hn-daily.db")
	err = db.Init(dbPath)
	if err != nil {
		fmt.Println("Something went wrong starting db connection:", err)
		os.Exit(1)
	}
	now := time.Now()
	posts, err := process.Day(&now)
	if err != nil {
		fmt.Println("Something went wrong processing today:", err)
		os.Exit(1)
	}
	j, err := json.Marshal(posts)
	if err != nil {
		fmt.Println("Something went wrong marshaling posts to JSON:", err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(fmt.Sprintf("HnPosts-%s.json", now.Format("2006-01-02")), j, 0644)
	if err != nil {
		fmt.Println("Something went wrong writing JSON to file:", err)
		os.Exit(1)
	}
}
