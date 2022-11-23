package accounts

import (
	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/Mirantis/dataeng/dataengctl/pkg/salesforce-reporting/accounts"
	"github.com/spf13/cobra"
)

const (
	ListExample = `
# this will print list of the boards
dataengctl salesforce accounts list

# this will create output in a yaml format
dataengctl salesforce accounts list --output yaml

`
)

func NewAccountsCommand(client *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "interact with accounts",
	}

	cmd.AddCommand(NewListCommand(client))
	return cmd
}

func NewListCommand(client *client.DataClient) *cobra.Command {
	options := &accounts.Options{
		Client: client,
	}
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List salesforce accounts",
		Example:       ListExample,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return options.ListAccounts()
		},
	}
	return cmd
}
