package client

import (
	//"database/sql"
	//"github.com/snowflakedb/gosnowflake"
	"github.com/Mirantis/dataeng/dataengctl/pkg/config"
)

type Option func(*Client)

func InjectSnowflakeAdaptorOption(adaptor ClientAdaptor) Option {
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
}

func newAdaptor(username, password, token string) ClientAdaptor {
	return nil
}

type adaptor struct {
}

var _ ClientAdaptor = &adaptor{}

type FakeAdaptor struct {
}

var _ ClientAdaptor = &FakeAdaptor{}
