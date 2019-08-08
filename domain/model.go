package domain

import "time"

type Post struct {
	Subreddit string
	Text      string
	CreatedAt time.Time
}

type User struct {
	Name string
}
