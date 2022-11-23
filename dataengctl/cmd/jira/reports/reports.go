package reports

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/Mirantis/dataeng/dataengctl/pkg/jira-reporting/reports"
)

const (
	ListExample = `
# this will print list of the boards
dataengctl jira boards list

# this will create output in a yaml format
dataengctl jira boards list --output yaml

`
)

func NewReportsCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "interact with reports",
	}

	cmd.AddCommand(NewListCommand(dataClient))
	cmd.AddCommand(NewExpandCommand(dataClient))
	return cmd
}

func NewListCommand(dataClient *client.DataClient) *cobra.Command {
	var output string
	var configPath string
	options := &reports.Options{
		Config: dataClient,
	}
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List jira reports",
		Example:       ListExample,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return options.List(os.Stdout)
		},
	}
	cmd.Flags().StringVarP(&configPath, "config-file", "c", "./config.yaml", "specify path to the config file, if left empty config.yaml from your current directory will be used")
	cmd.Flags().StringVarP(&output, "output", "o", "yaml", "use this flag to specify output")
	return cmd
}

func NewExpandCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "expand",
		Short:   "expand jira report",
		Example: "Some example how to run this",

		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("here is the list of expanded boards\n")
			return nil
		},
	}
	return cmd
}
