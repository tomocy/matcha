package client

import (
	"github.com/tomocy/matcha/app"
	"github.com/tomocy/matcha/infra"
)

func newPostUsecase() *app.PostUsecase {
	return &app.PostUsecase{
		Repo: infra.NewReddit(),
	}
}
