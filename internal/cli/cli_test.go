package cli

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func TestCli(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)))
	fmt.Println(client.Repositories.ListStatuses(ctx, "funapy-sandbox", "actions-sandbox", "ff9901ade9dcb344f7a5c94adaacbc5c865e7dd7", &github.ListOptions{}))
}
