package database

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	backupFileName     string
	backupFileNamePath string
)

var dbBackupCmd = &cobra.Command{
	Use:     "backup",
	Aliases: []string{},
	Short:   "Manage Database Backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	dbBackupCmd.AddCommand(dbBackupCreateCmd)

	dbBackupCreateCmd.Flags().StringVarP(&backupFileName, "name", "n", "", "back-up file name of database (default: <db-name>_backup.sql)")
	dbBackupCreateCmd.Flags().StringVarP(&dbType, "type", "t", "", "specify the type of database( Required)")
	dbBackupCreateCmd.Flags().StringVarP(&backupFileNamePath, "path", "P", ".", "path for the backup file (default: current directory)")
	dbBackupCreateCmd.Flags().Bool("no-data", false, "Exclude data from the backup (default: false)")
	dbBackupCreateCmd.Flags().Bool("no-create-info", false, "Exclude table schema from the backup (default: false)")
	dbBackupCreateCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "schedule backup of database")
	dbBackupCreateCmd.Flags().StringVarP(&dbUser, "user", "u", dbUser, "Database username")
	dbBackupCreateCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbBackupCreateCmd.MarkFlagsRequiredTogether("type", "user")
}
