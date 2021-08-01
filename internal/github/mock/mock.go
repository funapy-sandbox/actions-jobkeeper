package mock

import (
	"context"

	"github.com/funapy-sandbox/actions-jobkeeper/internal/github"
)

type Client struct {
	ListStatusesFunc func(ctx context.Context, owner, repo, ref string, opts *github.ListOptions) ([]*github.RepoStatus, *github.Response, error)
}

func (c *Client) ListStatuses(ctx context.Context, owner, repo, ref string, opts *github.ListOptions) ([]*github.RepoStatus, *github.Response, error) {
	return c.ListStatusesFunc(ctx, owner, repo, ref, opts)
}

var (
	_ github.Client = &Client{}
)
