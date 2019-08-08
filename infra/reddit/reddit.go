package reddit

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/tomocy/matcha/domain"
)

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
	SubredditNamePrefixed string        `json:"subreddit_name_prefixed"`
	AuthorFullname        string        `json:"author_fullname"`
	Title                 string        `json:"title"`
	CreatedUTC            unixTimestamp `json:"created_utc"`
}

func (p *Post) Adapt() *domain.Post {
	return &domain.Post{
		Subreddit: p.SubredditNamePrefixed,
		User: domain.User{
			Name: p.AuthorFullname,
		},
		Title:     p.Title,
		CreatedAt: time.Time(p.CreatedUTC),
	}
}

type unixTimestamp time.Time

func (t *unixTimestamp) UnmarshalJSON(data []byte) error {
	parsed, err := t.parse(string(data))
	if err != nil {
		return err
	}
	*t = unixTimestamp(parsed.Local())

	return nil
}

func (t *unixTimestamp) parse(ts string) (time.Time, error) {
	splited := strings.Split(ts, ".")
	if len(splited) != 2 {
		return time.Time{}, errors.New("invalid format of unix timestamp: the format should be sec.nsec")
	}
	sec, err := strconv.ParseInt(splited[0], 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	nsec, err := strconv.ParseInt(splited[1], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(sec, nsec), nil
}
