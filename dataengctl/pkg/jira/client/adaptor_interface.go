package client

import (
	"context"

	"github.com/andygrunwald/go-jira"
)


type AdaptorInterface interface {
	GetAllBoards(*jira.BoardListOptions) (*jira.BoardsList, *jira.Response, error)
	SearchIssuesWithContext(context.Context, string, *jira.SearchOptions) ([]jira.Issue, *jira.Response, error)
}
type adaptor struct {
	jiraClient *jira.Client
}

func (a *adaptor) SearchIssuesWithContext(ctx context.Context, jql string, options *jira.SearchOptions) ([]jira.Issue, *jira.Response, error) {
	return a.jiraClient.Issue.SearchWithContext(ctx, jql, options)
}

func (a *adaptor) GetAllBoards(opt *jira.BoardListOptions) (*jira.BoardsList, *jira.Response, error) {
	return a.jiraClient.Board.GetAllBoards(opt)
}
