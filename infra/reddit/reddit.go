package reddit

import "github.com/tomocy/matcha/domain"

type Posts struct {
	Data struct {
		Children []*struct {
			Data Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func (ps Posts) Adapt() []*domain.Post {
	adapteds := make([]*domain.Post, len(ps.Data.Children))
	for i, p := range ps.Data.Children {
		adapteds[i] = p.Data.Adapt()
	}

	return adapteds
}

type Post struct {
	SubredditNamePrefixed string `json:"subreddit_name_prefixed"`
	AuthorFullname        string `json:"author_fullname"`
	Title                 string `json:"title"`
}

func (p *Post) Adapt() *domain.Post {
	return &domain.Post{
		Title: p.Title,
	}
}
