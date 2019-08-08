package client

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tomocy/matcha/domain"
)

type CLI struct{}

func (c *CLI) PollPosts() error {
	u := newPostUsecase()
	ctx, cancelFn := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT)
	psCh, errCh := u.PollPosts(ctx)
	for {
		select {
		case ps := <-psCh:
			ordered := orderOlderPosts(ps)
			c.showPosts(ordered)
			fmt.Printf("updated at %s\n", time.Now().Format("2006/01/02 15:04"))
		case err := <-errCh:
			cancelFn()
			return err
		case sig := <-sigCh:
			cancelFn()
			fmt.Println(sig)
			return nil
		}
	}
}

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
