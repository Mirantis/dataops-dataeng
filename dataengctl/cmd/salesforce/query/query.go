package query

import (
	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/Mirantis/dataeng/dataengctl/pkg/salesforce-reporting/soql"
	"github.com/spf13/cobra"
)

const (
	ListExample = `
# this will print query output from SOQL string
dataengctl salesforce query "SOQL_string"

# this will query salesforce and show json output for SELECT+name+from+Account 
dataengctl salesforce query "SELECT+name+from+Account" 

`
)

func NewQueryCommand(client *client.DataClient) *cobra.Command {
	options := &soql.Options{
		Client: client,
	}
	cmd := &cobra.Command{
		Use:           "query",
		Short:         "List salesforce accounts",
		Example:       ListExample,
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return options.SOQLQuery(args[0])
		},
	}
	return cmd
}
