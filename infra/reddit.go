package infra

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/tomocy/deverr"

	"golang.org/x/oauth2"

	"github.com/tomocy/matcha/domain"
	"github.com/tomocy/matcha/infra/reddit"
)

func NewReddit() *Reddit {
	createWorkspace()
	return &Reddit{
		oauth: oauth{
			config: oauth2.Config{
				ClientID:     "w9IvCG5aiZb-fA",
				ClientSecret: "eYU7TZHoMmt1lOkOL9gbXBc2BTY",
				RedirectURL:  "http://localhost",
				Endpoint: oauth2.Endpoint{
					AuthURL:   "https://www.reddit.com/api/v1/authorize",
					TokenURL:  "https://www.reddit.com/api/v1/access_token",
					AuthStyle: oauth2.AuthStyleInHeader,
				},
				Scopes: []string{
					"read", "identity", "mysubreddits",
				},
			},
		},
	}
}

type Reddit struct {
	oauth oauth
}

type oauth struct {
	config oauth2.Config
	state  string
}

func (r *Reddit) FetchPosts() ([]*domain.Post, error) {
	return r.fetchPosts("https://oauth.reddit.com/new", nil)
}

func (r *Reddit) fetchPosts(destURL string, params url.Values) ([]*domain.Post, error) {
	tok, err := r.trieveAuthorization()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %s", err)
	}

	var posts reddit.Posts
	if err := r.trieve(&oauthRequest{
		tok:     tok,
		method:  http.MethodGet,
		destURL: destURL,
		params:  params,
	}, &posts); err != nil {
		r.resetConfig()
		return nil, fmt.Errorf("failed to fetch posts: %s", err)
	}

	if err := r.saveConfig(redditConfig{
		AccessToken: tok,
	}); err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %s", err)
	}

	return posts.Adapt(), nil
}

func (r *Reddit) trieveAuthorization() (*oauth2.Token, error) {
	if config, err := loadConfig(); err == nil {
		return config.Reddit.AccessToken, nil
	}

	r.oauth.state = fmt.Sprintf("%d", rand.Intn(1000))
	url := r.oauth.config.AuthCodeURL(r.oauth.state)
	fmt.Printf("open this link: %s\n", url)
	tokCh, errCh := r.handleAuthorizationRedirect()
	select {
	case tok := <-tokCh:
		return tok, nil
	case err := <-errCh:
		return nil, err
	}
}

func (r *Reddit) handleAuthorizationRedirect() (<-chan *oauth2.Token, <-chan error) {
	tokCh, errCh := make(chan *oauth2.Token), make(chan error)
	go func() {
		defer func() {
			close(tokCh)
			close(errCh)
		}()

		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			q := req.URL.Query()
			state, code := q.Get("state"), q.Get("code")
			if r.oauth.state != state {
				w.WriteHeader(http.StatusBadRequest)
				errCh <- errors.New("invalid state")
				return
			}
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
				Transport: new(oauthUserAgentTransport),
			})
			tok, err := r.oauth.config.Exchange(ctx, code)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errCh <- err
				return
			}

			tokCh <- tok
		})
		if err := http.ListenAndServe(":80", nil); err != nil {
			errCh <- err
		}
	}()

	return tokCh, errCh
}

func (r *Reddit) trieve(req *oauthRequest, dest interface{}) error {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
		Transport: new(oauthUserAgentTransport),
	})
	client := r.oauth.config.Client(ctx, req.tok)
	var resp *http.Response
	var err error
	switch req.method {
	case http.MethodGet:
		resp, err = client.Get(req.destURL)
	default:
		return deverr.NotImplemented
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return readJSON(resp.Body, dest)
}

func (r *Reddit) saveConfig(conf redditConfig) error {
	if loaded, err := loadConfig(); err == nil {
		loaded.Reddit = conf
		return saveConfig(loaded)
	}

	return saveConfig(&config{
		Reddit: conf,
	})
}

func (r *Reddit) resetConfig() error {
	conf := new(config)
	if loaded, err := loadConfig(); err == nil {
		loaded.Reddit = redditConfig{}
		conf = loaded
	}

	return saveConfig(conf)
}

type oauthRequest struct {
	tok             *oauth2.Token
	method, destURL string
	params          url.Values
}

type oauthUserAgentTransport struct{}

func (t *oauthUserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "oauth-client/0.0")

	return http.DefaultTransport.RoundTrip(r)
}
