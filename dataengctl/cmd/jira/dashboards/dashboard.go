package dashboards

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/Mirantis/dataeng/dataengctl/pkg/jira-reporting/dashboards"
)

const (
	ListExample = `
# this will print list of the dashboard
dataengctl jira dashboard list

# this will create output in a yaml format
dataengctl jira dashboard list --output yaml

`
)

func NewDashboardCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "interact with dashboard",
	}

	cmd.AddCommand(NewListCommand(dataClient))
	cmd.AddCommand(NewExpandCommand(dataClient))
	return cmd
}

func NewListCommand(dataClient *client.DataClient) *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List jira dashboard",
		Example:       ListExample,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return dashboard.List(dataClient)
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "yaml", "use this flag to specify output")
	return cmd
}

func NewExpandCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "expand",
		Short:   "expand jira dashboard",
		Example: "Some example how to run this",

		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("here is the list of expanded dashboards\n")
			return nil
		},
	}
	return cmd
}
