package client

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Mirantis/dataeng/dataengctl/pkg/config"
	"github.com/andygrunwald/go-jira"
)

type Option func(*Client)

func InjectJiraAdaptorOption(adaptor AdaptorInterface) Option {
	return func(c *Client) {
		c.adaptor = adaptor
	}
}

func NewClient(cfg *config.JiraConfig, opts ...Option) (ClientInterface, error) {
	client := &Client{}

	for _, o := range opts {
		o(client)
	}

	if client.adaptor == nil {
		tp := jira.BasicAuthTransport{
			Username: cfg.Username,
			Password: cfg.Token,
		}
		jiraClient, err := jira.NewClient(tp.Client(), cfg.URL)
		if err != nil {
			return nil, err
		}
		client.adaptor = &adaptor{jiraClient: jiraClient}
	}

	return client, nil
}

type Client struct {
	Config  *config.JiraConfig
	adaptor AdaptorInterface
}

type ClientInterface interface {
	IssueIterator(opts IssueIteratorOptions) (IssueIteratorInterface, error)
	GetAllBoards(opt *jira.BoardListOptions) (*jira.BoardsList, *jira.Response, error)
	DashboardBytes() ([]byte, error)
	ProjectBytes() ([]byte, error)
}

type IssueIteratorInterface interface {
	Next(ctx context.Context) (*jira.Issue, error)
}

func (c *Client) IssueIterator(opts IssueIteratorOptions) (IssueIteratorInterface, error) {
	return NewIssueIterator(context.Background(), c.adaptor, opts)
}

func (c *Client) GetAllBoards(opt *jira.BoardListOptions) (*jira.BoardsList, *jira.Response, error) {
	return c.adaptor.GetAllBoards(&jira.BoardListOptions{})
}

func (c *Client) DashboardBytes() ([]byte, error) {
	url := strings.Join([]string{c.Config.URL, "rest/api/latest/dashboard"}, "/")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Config.Username, c.Config.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c *Client) ProjectBytes() ([]byte, error) {
	url := strings.Join([]string{c.Config.URL, "rest/api/latest/project"}, "/")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Config.Username, c.Config.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
