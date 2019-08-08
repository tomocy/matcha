package domain

import "time"

type Post struct {
	Text      string
	CreatedAt time.Time
}

type User struct {
	Name string
}
