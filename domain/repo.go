package domain

import "context"

type PostRepo interface {
	PolePosts(context.Context) (<-chan []*Post, <-chan error)
	FetchPosts() ([]*Post, error)
}
