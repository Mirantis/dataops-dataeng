package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIteratorNext(t *testing.T) {
	testCases := []struct {
		name          string
		expectedError string

		adaptor AdaptorInterface
		options IssueIteratorOptions
	}{
		{
			name:          "error passhrough",
			expectedError: "next error",

			options: IssueIteratorOptions{
				PageSize: 3,
				Query:    "something",
				Expand:   "changelog",
			},
			adaptor: &FakeAdaptor{
				searchIssuesWithContextCall: func(ctx context.Context, jql string, options *jira.SearchOptions) ([]jira.Issue, *jira.Response, error) {
					return nil, nil, fmt.Errorf("next error")
				},
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			iterator, err := NewIssueIterator(context.Background(), tt.adaptor, tt.options)
			require.NoError(t, err)
			require.NotNil(t, iterator)
			for {
				issue, actualErr := iterator.Next(context.TODO())
				if tt.expectedError != "" {
					require.Error(t, actualErr)
					assert.Contains(t, actualErr.Error(), tt.expectedError)
					break
				} else {
					require.NoError(t, actualErr)
					require.NotNil(t, issue)
					if issue == nil {
						break
					}
				}

			}

		})
	}

}
