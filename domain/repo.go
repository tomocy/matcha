package domain

import "context"

type PostRepo interface {
	PollPosts(context.Context) (<-chan []*Post, <-chan error)
	FetchPosts() ([]*Post, error)
}
