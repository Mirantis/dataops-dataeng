package client

import (
	"github.com/Mirantis/dataeng/dataengctl/pkg/config"
	"github.com/simpleforce/simpleforce"
)

type Option func(*Client)

func InjectSalesforceAdaptorOption(adaptor ClientAdaptor) Option {
	return func(c *Client) {
		c.adaptor = adaptor
	}
}

type Client struct {
	Config  *config.SalesForceConfig
	adaptor ClientAdaptor
}

func NewClient(config *config.SalesForceConfig, opts ...Option) *Client {
	client := &Client{
		Config: config,
	}
	for _, o := range opts {
		o(client)
	}
	if client.adaptor == nil {
		client.adaptor = newAdaptor(config.URL, config.ClientID, config.APIVersion)
	}
	return client
}

type ClientAdaptor interface {
	LoginPassword(string, string, string) error
	Query(string) (*simpleforce.QueryResult, error)
}

func (c *Client) Query(q string) (*simpleforce.QueryResult, error) {
	if err := c.adaptor.LoginPassword(c.Config.Username, c.Config.Password, c.Config.Token); err != nil {
		return nil, err
	}

	// q := "SELECT+Environment2__r.Name+FROM+Case"
	return c.adaptor.Query(q) // Note: for Tooling API, use client.Tooling().Query(q)
}

func newAdaptor(username, password, token string) ClientAdaptor {
	return &adaptor{
		Client: simpleforce.NewClient(username, password, token),
	}
}

type adaptor struct {
	*simpleforce.Client
}

var _ ClientAdaptor = &adaptor{}

type FakeAdaptor struct {
	loginPasswordCall func(string, string, string) error
	queryCall         func(string) (*simpleforce.QueryResult, error)
}

var _ ClientAdaptor = &FakeAdaptor{}

func (a *FakeAdaptor) LoginPassword(url, clientID, token string) error {
	return a.loginPasswordCall(url, clientID, token)
}

func (a *FakeAdaptor) Query(q string) (*simpleforce.QueryResult, error) {
	return a.queryCall(q)
}
