package database

import (
	"errors"

	"github.com/spf13/cobra"
)

var dbLogCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{},
	Short:   "Manage Logs related to backup and restoration process",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	dbLogCmd.AddCommand(dbLogsListCmd)
	dbLogCmd.AddCommand(dbLogListDeletionCmd)

	dbLogsListCmd.Flags().StringVarP(&dbType, "type", "t", "", "specify the type of database [*Required]")
	dbLogsListCmd.Flags().StringVarP(&dbUser, "user", "u", "", "Database username [*Required]")
	dbLogsListCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbLogListDeletionCmd.Flags().StringVarP(&dbType, "type", "t", "", "specify the type of database [*Required]")
	dbLogListDeletionCmd.Flags().StringVarP(&dbUser, "user", "u", "", "Database username [*Required]")
	dbLogListDeletionCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbLogListDeletionCmd.MarkFlagsRequiredTogether("type", "user")
	dbLogsListCmd.MarkFlagsRequiredTogether("type", "user")
}
