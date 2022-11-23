package client

import (
	"context"

	"github.com/andygrunwald/go-jira"
)

type FakeAdaptor struct {
	searchIssuesWithContextCall func(ctx context.Context, jql string, options *jira.SearchOptions) ([]jira.Issue, *jira.Response, error)
	getAllBoardsCall            func(opt *jira.BoardListOptions) (*jira.BoardsList, *jira.Response, error)
}

func (a *FakeAdaptor) SearchIssuesWithContext(ctx context.Context, jql string, options *jira.SearchOptions) ([]jira.Issue, *jira.Response, error) {
	return a.searchIssuesWithContextCall(ctx, jql, options)
}

func (a *FakeAdaptor) GetAllBoards(opt *jira.BoardListOptions) (*jira.BoardsList, *jira.Response, error) {
	return a.getAllBoardsCall(opt)
}
