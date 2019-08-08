package app

import (
	"context"

	"github.com/tomocy/matcha/domain"
)

type PostUsecase struct {
	Repo domain.PostRepo
}

func (u *PostUsecase) PollPosts(ctx context.Context) (<-chan []*domain.Post, <-chan error) {
	return u.Repo.PollPosts(ctx)
}

func (u *PostUsecase) FetchPosts() ([]*domain.Post, error) {
	return u.Repo.FetchPosts()
}
