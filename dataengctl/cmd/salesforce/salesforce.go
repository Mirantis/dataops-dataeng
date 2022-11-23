package salesforce

import (
	"github.com/Mirantis/dataeng/dataengctl/cmd/salesforce/accounts"
	"github.com/Mirantis/dataeng/dataengctl/cmd/salesforce/query"
	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/spf13/cobra"
)

// Wrapper for Jira Client
func NewSalesforceCmd(client *client.DataClient) *cobra.Command {
	//var configFile string
	//sf_cfg := &config.Config{}
	cmd := &cobra.Command{
		Use:   "salesforce",
		Short: "Interact with salesforce",
		//PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//	return config.ReadConfig(configFile, sf_cfg)
		//},
	}

	// Wrapper for Output Flag to All Base Commands
	/*func OutPutCmd() *cobra.Command {
		var output string
		output := &output{}
		cmd := &cobra.Command{
			Use:	"output"
		},
	}*/

	// Universal Flag Commands
	//cmd.PersistentFlags().StringVar(&configFile, "config-file", "", "path to config file")
	//cmd.PersistentFlags().StringVar(&output), "output", "", "csv,dataframe,json,yaml")

	// Universal Base Commands
	//cmd.AddCommand(boards.NewBoardsCommand(cfg))
	cmd.AddCommand(accounts.NewAccountsCommand(client))
	cmd.AddCommand(query.NewQueryCommand(client))
	//cmd.AddCommand(dashboards.NewDashboardCommand(cfg))
	//cmd.AddCommand(issues.NewIssueCommand(cfg))
	//cmd.AddCommand(login.NewLoginCommand(cfg))
	//cmd.AddCommand(projects.NewProjectCommand(cfg))
	//cmd.AddCommand(reports.NewReportsCommand(cfg))
	//cmd.AddCommand(epics.NewEpicCommand())
	return cmd
}
