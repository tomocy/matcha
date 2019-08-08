package client

import (
	"fmt"

	"github.com/tomocy/matcha/domain"
)

type CLI struct{}

func (c *CLI) FetchPosts() error {
	u := newPostUsecase()
	ps, err := u.FetchPosts()
	if err != nil {
		return err
	}

	ordered := orderOlderPosts(ps)
	c.showPosts(ordered)

	return nil
}

func (c *CLI) showPosts(ps []*domain.Post) {
	strable := asciiPosts(ps)
	fmt.Print(strable)
}
