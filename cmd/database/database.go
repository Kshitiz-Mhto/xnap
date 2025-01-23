package database

import (
	"errors"

	"github.com/spf13/cobra"
)

// DBCmd is the root command for the db subcommand
var DBCmd = &cobra.Command{
	Use:     "database",
	Aliases: []string{"db", "databases"},
	Short:   "Manage Databases",
	Long:    `operation related to databases`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {

	DBCmd.AddCommand(dbListCmd)

	dbListCmd.Flags().StringVarP(&dbType, "type", "t", "all", "Filter by database type (all/mysql/postgres)")

}
