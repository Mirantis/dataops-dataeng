package client

import (
	"context"
	"fmt"

	"github.com/andygrunwald/go-jira"
)

func NewIssueIterator(ctx context.Context, adp AdaptorInterface, opts IssueIteratorOptions) (IssueIteratorInterface, error) {
	if opts.Query == "" {
		return nil, fmt.Errorf("must pass non-nil query")
	}

	it := &IssueIterator{
		adaptor:  adp,
		query:    opts.Query,
		expand:   opts.Expand,
		pageSize: opts.PageSize,
	}
	return it, it.more(ctx)
}

type IssueIterator struct {
	adaptor      AdaptorInterface
	query        string
	expand       string
	issues       []jira.Issue
	pageSize     int
	startAt      int
	endAt        int
	totalResults int
	index        int
}

type IssueIteratorOptions struct {
	Query    string
	Expand   string
	PageSize int
}

func (it *IssueIterator) initialize() error {
	return nil
}

func (it *IssueIterator) more(ctx context.Context) error {
	issues, response, err := it.adaptor.SearchIssuesWithContext(
		ctx,
		it.query,
		&jira.SearchOptions{
			StartAt:    it.endAt,
			MaxResults: it.pageSize,
			Expand:     it.expand,
		},
	)
	if err != nil {
		return err
	}

	it.issues = issues
	it.startAt = response.StartAt
	it.totalResults = response.Total
	it.endAt = it.endAt + response.MaxResults
	if it.endAt > response.Total {
		it.endAt = response.Total
	}
	return nil
}

func (it *IssueIterator) Next(ctx context.Context) (*jira.Issue, error) {
	it.index = it.index + 1

	if it.index > it.totalResults {
		return nil, nil
	}

	if it.index > it.endAt {
		if err := it.more(ctx); err != nil {
			return nil, fmt.Errorf("error getting issues: %v", err)
		}
	}
	return &it.issues[it.index-it.startAt-1], nil
}
