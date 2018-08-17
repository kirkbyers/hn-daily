package db

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

var (
	tmpDbPath   string
	testHnScrap []HnPost
)

func TestSaveHnScrape(t *testing.T) {
	setup()
	defer teardown()
	err := SaveHnScrape(testHnScrap)
	if err != nil {
		t.Error(err)
	}
	now := time.Now()
	actual, err := GetScrapesForDay(&now)
	expect := testHnScrap
	if err != nil {
		t.Error(err)
	}
	if len(expect) != len(actual) {
		t.Errorf("Len of expected is %d and len of actual is %d", len(expect), len(actual))
		t.Failed()
	}
	for i, v := range actual {
		e := expect[i]
		if v != e {
			t.Failed()
		}
	}
}

func setup() {
	rand.Seed(time.Now().UnixNano())
	// Set-up tmp db
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Something went wrong finding homeDir:", err)
		os.Exit(1)
	}
	tmpDbPath = filepath.Join(home, "tmp.db")
	err = Init(tmpDbPath)
	if err != nil {
		fmt.Println("Something went wrong creating db:", err)
		os.Exit(1)
	}
	testHnScrap = generateMockScrape()
}

func generateMockScrape() (result []HnPost) {
	// Set mock data
	for i := 0; i < 30; i++ {
		result = append(result, HnPost{
			Title: fmt.Sprintf("Post #%d", i),
			Score: rand.Intn(9999),
			ID:    string(i),
		})
	}
	return result
}

func teardown() {
	err := os.Remove(tmpDbPath)
	if err != nil {
		fmt.Println("There was a problem deleting tmp db:", err)
		os.Exit(1)
	}
}
