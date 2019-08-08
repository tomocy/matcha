package domain

import "time"

type Post struct {
	Subreddit string
	User      User
	Text      string
	CreatedAt time.Time
}

type User struct {
	Name string
}
