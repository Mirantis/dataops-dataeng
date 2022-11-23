package boards

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/Mirantis/dataeng/dataengctl/pkg/jira-reporting/boards"
)

const (
	ListExample = `
# this will print list of the boards
dataengctl jira boards list

# this will create output in a yaml format
dataengctl jira boards list --output yaml

`
)

func NewBoardsCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "board",
		Short: "interact with boards",
	}

	cmd.AddCommand(NewListCommand(dataClient))
	cmd.AddCommand(NewExpandCommand(dataClient))
	return cmd
}

func NewListCommand(dataClient *client.DataClient) *cobra.Command {
	var output string
	options := &boards.Options{
		Config: dataClient,
	}
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List jira boards",
		Example:       ListExample,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return options.List(os.Stdout)
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "yaml", "use this flag to specify output")
	return cmd
}

func NewExpandCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "expand",
		Short:   "expand jira board",
		Example: "Some example how to run this",

		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("here is the list of expanded boards\n")
			return nil
		},
	}
	return cmd
}
