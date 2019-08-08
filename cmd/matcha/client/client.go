package client

import (
	"fmt"
	"strings"

	"github.com/tomocy/matcha/app"
	"github.com/tomocy/matcha/domain"
	"github.com/tomocy/matcha/infra"
	"github.com/tomocy/tago"
)

func newPostUsecase() *app.PostUsecase {
	return &app.PostUsecase{
		Repo: infra.NewReddit(),
	}
}

type asciiPosts []*domain.Post

func (ps asciiPosts) String() string {
	tago := tago.NewWithout(tago.DefaultDuration, "2006/01/02 15:04")
	var b strings.Builder
	for i, p := range ps {
		if i == 0 {
			b.WriteString("----------\n")
		}
		b.WriteString(fmt.Sprintf("%s %s\n%s\n%s\n", p.Subreddit, tago.Ago(p.CreatedAt), p.User.Name, p.Title))
		b.WriteString("----------\n")
	}

	return b.String()
}
