package client

import (
	"fmt"
	"testing"

	"github.com/Mirantis/dataeng/dataengctl/pkg/config"

	"github.com/simpleforce/simpleforce"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {

	testCases := []struct {
		name          string
		query         string
		expectedError string

		config  *config.SalesForceConfig
		adaptor ClientAdaptor
	}{
		{
			name: "success credentials are passed to LoginPassword() call",
			config: &config.SalesForceConfig{
				Password: "my-pass",
				Token:    "my-token",
				Username: "my-username",
			},
			query: "something",
			adaptor: &FakeAdaptor{
				loginPasswordCall: func(username, pass, token string) error {
					if username != "my-username" || pass != "my-pass" || token != "my-token" {
						return fmt.Errorf("credentials are not passed")
					}
					return nil
				},
				queryCall: func(s string) (*simpleforce.QueryResult, error) {
					return &simpleforce.QueryResult{}, nil
				},
			},
		},
		{
			name: "error LoginPassword()",
			config: &config.SalesForceConfig{
				Password: "my-pass",
				Token:    "my-token",
				Username: "my-username",
			},
			query:         "something",
			expectedError: "LoginPassword error",
			adaptor: &FakeAdaptor{
				loginPasswordCall: func(username, pass, token string) error {
					return fmt.Errorf("LoginPassword error")
				},
			},
		},
		{
			name: "error Query()",
			config: &config.SalesForceConfig{
				Password: "my-pass",
				Token:    "my-token",
				Username: "my-username",
			},
			query:         "something",
			expectedError: "Query error",
			adaptor: &FakeAdaptor{
				loginPasswordCall: func(username, pass, token string) error {
					return nil
				},
				queryCall: func(s string) (*simpleforce.QueryResult, error) {
					return nil, fmt.Errorf("Query error")
				},
			},
		},
		{
			name: "error default adaptor wrong protocol scheme",
			config: &config.SalesForceConfig{
				URL:      "nonexistant new url",
				Password: "my-pass",
				Token:    "my-token",
				Username: "my-username",
			},
			query:         "something",
			expectedError: "unsupported protocol scheme",
		},
	}

	for _, tt := range testCases {
		tt := tt
		client := NewClient(tt.config, InjectSalesforceAdaptorOption(tt.adaptor))

		result, err := client.Query(tt.query)
		if tt.expectedError != "" {
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, result)
		}

	}
}
