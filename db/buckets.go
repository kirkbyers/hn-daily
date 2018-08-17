package db

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

var (
	dailyBuckets = []byte("daily")
	db           *bolt.DB
)

type dailyRecord []HnPost

// HnPost is what is scraped per post from hn
type HnPost struct {
	Title string
	Score int
	URL   string
	ID    string
}

// Init Creates an instance of boltDB
func Init(dbFilePath string) (err error) {
	db, err = bolt.Open(dbFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
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
		now := time.Now().Format(time.RFC3339)
		key := []byte(now)
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
