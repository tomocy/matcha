package client

import (
	"strings"

	"github.com/tomocy/matcha/app"
	"github.com/tomocy/matcha/domain"
	"github.com/tomocy/matcha/infra"
)

func newPostUsecase() *app.PostUsecase {
	return &app.PostUsecase{
		Repo: infra.NewReddit(),
	}
}

type asciiPosts []*domain.Post

func (ps asciiPosts) String() string {
	var b strings.Builder
	for i, p := range ps {
		if i == 0 {
			b.WriteString("-----")
		}
		b.WriteString(p.Text)
		b.WriteString("-----")
	}

	return b.String()
}
