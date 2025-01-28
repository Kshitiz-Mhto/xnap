package database

import (
	"errors"

	"github.com/Kshitiz-Mhto/dsync/pkg/config"
	"github.com/spf13/cobra"
)

var (
	MySQL_DB_USER     string = config.Envs.MySQL_DB_USER
	MySQL_DB_PASSWORD string = config.Envs.MySQL_DB_PASSWORD
	MySQL_DB_HOST     string = config.Envs.MySQL_DB_HOST
	MySQL_DB_PORT     string = config.Envs.MySQL_DB_PORT

	POSTGRES_DB_USER     string = config.Envs.POSTGRES_DB_USER
	POSTGRES_DB_PASSWORD string = config.Envs.POSTGRES_DB_PASSWORD
	POSTGRES_DB_HOSTOST  string = config.Envs.POSTGRES_DB_HOST
	POSTGRES_DB_PORT     string = config.Envs.POSTGRES_DB_PORT
)

var (
	dbName   string
	dbOwner  string
	dbType   string
	schedule string
	masterDB string
	slaveDB  string
)

// DBCmd is the root command for the db subcommand
var DBCmd = &cobra.Command{
	Use:     "database",
	Aliases: []string{"db", "databases"},
	Short:   "Manage Databases",
	Long:    `Create, List and Delete databases`,
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
	DBCmd.AddCommand(dbCreateCmd)
	DBCmd.AddCommand(dbDeleteCmd)
	DBCmd.AddCommand(dbBackupCmd)
	DBCmd.AddCommand(dbRestoreCmd)
	DBCmd.AddCommand(dbSyncCmd)

	dbListCmd.Flags().StringVarP(&dbType, "type", "t", "all", "Filter by database type (all/mysql/postgres)")

	dbCreateCmd.Flags().StringVarP(&dbType, "type", "t", "mysql", "Create database type MySQL")
	dbCreateCmd.Flags().StringVarP(&dbOwner, "owner", "o", "postgres", "Specify owner only for postgres database")

	dbDeleteCmd.Flags().StringVarP(&dbType, "type", "t", "", "Specify the database type for deletion")

	dbRestoreCmd.Flags().StringVarP(&backupFullFilePath, "backup", "b", "", "Path to the backup file (required)")
	dbRestoreCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "Time to schedule the restoration (in a format like HH:MM or a cron-like string)")
	dbRestoreCmd.Flags().StringVarP(&dbType, "type", "t", "", "Specify the database type for restoration")

	dbSyncCmd.Flags().StringVarP(&masterDB, "master", "m", "", "Specify the source/master database")
	dbSyncCmd.Flags().StringVarP(&slaveDB, "slave", "s", "", "Specify the destination/slave database")
	dbSyncCmd.Flags().StringVarP(&dbType, "type", "", "", "Specify the database type")

	dbDeleteCmd.MarkFlagRequired("type")

	dbRestoreCmd.MarkFlagRequired("backup")
	dbRestoreCmd.MarkFlagRequired("type")

	dbSyncCmd.MarkFlagsRequiredTogether("master", "slave", "type")

}
