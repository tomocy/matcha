package domain

import "time"

type Post struct {
	Subreddit string
	User      User
	Title     string
	CreatedAt time.Time
}

type User struct {
	Name string
}
