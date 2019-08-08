package app

import "github.com/tomocy/matcha/domain"

type PostUsecase struct {
	repo domain.PostRepo
}

func (u *PostUsecase) FetchPosts() ([]*domain.Post, error) {
	return u.repo.FetchPosts()
}
