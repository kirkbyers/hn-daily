package process

import (
	"sort"
	"time"

	"github.com/kirkbyers/hn-daily/db"
)

// Day returns the top 30 post 24 hr from given time
func Day(t *time.Time) (result []db.HnPost, err error) {
	var posts []db.HnPost
	posts, err = db.GetScrapesForDay(t)

	// Build a map of post with highest seen scores of the day
	postMap := map[string]db.HnPost{}
	for _, v := range posts {
		if post, ok := postMap[v.ID]; ok {
			old := &post
			new := &v
			if new.Score > old.Score {
				postMap[v.ID] = *new
			}
			continue
		}
		postMap[v.ID] = v
	}

	// Build sortedPosts array
	var sortedPosts []db.HnPost
	for _, post := range postMap {
		sortedPosts = append(sortedPosts, post)
	}
	// Sort Sorted Posts
	sort.Slice(sortedPosts, func(i, j int) bool {
		return sortedPosts[i].Score > sortedPosts[j].Score
	})

	// Cut sorted posts to length of normal Frontpage
	if len(sortedPosts) >= 30 {
		result = sortedPosts[:30]
	} else {
		result = sortedPosts
	}
	return result, nil
}
