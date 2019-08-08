package infra

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

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
				RedirectURL:  "http://localhost/reddit/authorization",
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

func (r *Reddit) PolePosts(ctx context.Context) (<-chan []*domain.Post, <-chan error) {
	return r.polePosts(ctx, "https://oauth.reddit.com/new")
}

func (r *Reddit) polePosts(ctx context.Context, dest string) (<-chan []*domain.Post, <-chan error) {
	psCh, errCh := make(chan []*domain.Post), make(chan error)
	go func() {
		defer func() {
			close(psCh)
			close(errCh)
		}()

		sendPosts := func(lastID string, psCh chan<- []*domain.Post, errCh chan<- error) string {
			params := make(url.Values)
			if lastID != "" {
				params.Set("after", lastID)
			}
			ps, err := r.fetchPosts(dest, params)
			if err != nil {
				errCh <- err
				return ""
			}
			if len(ps) <= 0 {
				return ""
			}

			psCh <- ps
			return ps[0].ID
		}

		lastID := sendPosts("", psCh, errCh)
		for {
			select {
			case <-ctx.Done():
				break
			case <-time.After(2 * time.Minute):
				lastID = sendPosts(lastID, psCh, errCh)
			}
		}
	}()

	return psCh, errCh
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
	if config, err := loadConfig(); err == nil && config.Reddit.AccessToken.Valid() {
		return config.Reddit.AccessToken, nil
	}

	r.oauth.state = fmt.Sprintf("%d", rand.Intn(1000))
	url := r.oauth.config.AuthCodeURL(r.oauth.state, oauth2.SetAuthURLParam("duration", "permanent"))
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

		http.HandleFunc("/reddit/authorization", func(w http.ResponseWriter, req *http.Request) {
			q := req.URL.Query()
			state, code := q.Get("state"), q.Get("code")
			if r.oauth.state != state {
				w.WriteHeader(http.StatusBadRequest)
				errCh <- errors.New("invalid state")
				return
			}
			ctx := r.userAgentTransportContext(context.Background())
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
	resp, err := r.request(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if http.StatusBadRequest <= resp.StatusCode {
		return errors.New(resp.Status)
	}

	return readJSON(resp.Body, dest)
}

func (r *Reddit) request(req *oauthRequest) (*http.Response, error) {
	ctx := r.userAgentTransportContext(context.Background())
	client := r.oauth.config.Client(ctx, req.tok)
	if req.method != http.MethodGet {
		return client.PostForm(req.destURL, req.params)
	}

	parsed, err := url.Parse(req.destURL)
	if err != nil {
		return nil, err
	}
	if req.params != nil {
		parsed.RawQuery = req.params.Encode()
	}

	return client.Get(parsed.String())
}

func (r *Reddit) userAgentTransportContext(parent context.Context) context.Context {
	return context.WithValue(parent, oauth2.HTTPClient, &http.Client{
		Transport: new(oauthUserAgentTransport),
	})
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
