package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type (
	ListOptions = github.ListOptions
	RepoStatus  = github.RepoStatus
	Response    = github.Response
)

type Client interface {
	ListStatuses(ctx context.Context, owner, repo, ref string, opts *ListOptions) ([]*RepoStatus, *Response, error)
}

type client struct {
	ghc *github.Client
}

func NewClient(ctx context.Context, token string) Client {
	return &client{
		ghc: github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		))),
	}
}

func (c *client) ListStatuses(ctx context.Context, owner, repo, ref string, opts *ListOptions) ([]*RepoStatus, *Response, error) {
	return c.ghc.Repositories.ListStatuses(ctx, owner, repo, ref, opts)
}
