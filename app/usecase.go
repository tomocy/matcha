package app

import "github.com/tomocy/matcha/domain"

type PostUsecase struct {
	Repo domain.PostRepo
}

func (u *PostUsecase) FetchPosts() ([]*domain.Post, error) {
	return u.Repo.FetchPosts()
}
