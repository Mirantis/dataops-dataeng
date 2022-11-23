package client

import (
	// "context"
	// "golang.org/x/oauth2"

	"github.com/Mirantis/dataeng/dataengctl/pkg/config"
)

type Option func(*Client)

func InjectGitHubAdaptorOption(adaptor ClientAdaptor) Option {
	return func(c *Client) {
		c.adaptor = adaptor
	}
}

type Client struct {
	Config  *config.GitHubConfig
	adaptor ClientAdaptor
}

func NewClient(config *config.GitHubConfig, opts ...Option) *Client {
	client := &Client{
		Config: config,
	}
	for _, o := range opts {
		o(client)
	}
	if client.adaptor == nil {
		client.adaptor = newAdaptor(config.OrgID, config.Token)
	}
	return client
}

type ClientAdaptor interface {
}

func newAdaptor(orgID, token string) ClientAdaptor {
	return nil
}

type adaptor struct {
}

var _ ClientAdaptor = &adaptor{}

type FakeAdaptor struct {
}

var _ ClientAdaptor = &FakeAdaptor{}

// func OAuthClient() {
// 	ctx := context.Background()
// 	ts := oauth2.StaticTokenSource(
// 		&oauth2.Token{AccessToken: "... your access token ..."},
// 	)
// 	tc := oauth2.NewClient(ctx, ts)
// 	client := github.NewClient(tc)
// 	// list all repositories for the authenticated user
// }
