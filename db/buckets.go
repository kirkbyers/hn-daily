package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	homedir "github.com/mitchellh/go-homedir"
)

var dailyBuckets = []byte("daily")
var db *bolt.DB

type dailyRecord []HnPost

// HnPost is what is scraped per post from hn
type HnPost struct {
	Title string
	Score int
	URL   string
	ID    string
}

// Init Creates an instance of boltDB
func Init() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, "hn-daily.db")
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(dailyBuckets)
		return err
	})
}

// SaveHnScrape saves scraped data to boltDB
func SaveHnScrape(hnScrape []HnPost) error {
	err := db.Update(func(tx *bolt.Tx) error {
		buc := tx.Bucket(dailyBuckets)
		// Set key for bucket the current time formated to a sortable string
		key := []byte(time.Now().Format(time.RFC3339))
		// Create JSON encoder
		buf, err := json.Marshal(hnScrape)
		if err != nil {
			return err
		}
		return buc.Put(key, buf)
	})
	return err
}

// GetScrapesForDay gets all the HN posts for 1 day
func GetScrapesForDay(t *time.Time) (result []HnPost, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(dailyBuckets).Cursor()
		min := []byte(t.Add(-24 * time.Hour).Format(time.RFC3339))
		max := []byte(t.Format(time.RFC3339))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			fmt.Printf("%s: %s\n", k, v)
			var posts []HnPost
			d := json.NewDecoder(bytes.NewBuffer(v))
			if err := d.Decode(&posts); err != nil {
				return err
			}
			result = append(result, posts...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
