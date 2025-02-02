package database

import (
	"errors"
	"time"

	"github.com/Kshitiz-Mhto/xnap/pkg/config"
	"github.com/spf13/cobra"
)

var (
	MySQL_DB_USER     string = config.Envs.MySQL_DB_USER
	MySQL_DB_PASSWORD string = config.Envs.MySQL_DB_PASSWORD
	MySQL_DB_HOST     string = config.Envs.MySQL_DB_HOST
	MySQL_DB_PORT     string = config.Envs.MySQL_DB_PORT

	POSTGRES_DB_USER     string = config.Envs.POSTGRES_DB_USER
	POSTGRES_DB_PASSWORD string = config.Envs.POSTGRES_DB_PASSWORD
	POSTGRES_DB_HOST     string = config.Envs.POSTGRES_DB_HOST
	POSTGRES_DB_PORT     string = config.Envs.POSTGRES_DB_PORT
)

var (
	dbName       string
	dbOwner      string
	dbType       string
	dbUser       string
	dbPassword   string
	promptPass   bool
	schedule     string
	status       string
	errorMessage string
	command      string
	start        time.Time
	duration     float64
)

// DBCmd is the root command for the db subcommand
var DBCmd = &cobra.Command{
	Use:     "database",
	Aliases: []string{"db", "databases"},
	Short:   "Manage Databases",
	Long:    `Create, List, Delete databases And backup and restore your database with config only or data only or whole and create dump file also provide feature of scheduling.`,
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

	dbListCmd.Flags().StringVarP(&dbType, "type", "t", "", "Specify type of database type: `mysql` or `postgres`. [*Required]")
	dbListCmd.Flags().StringVarP(&dbUser, "user", "u", "", "Database username [*Required]")
	dbListCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbCreateCmd.Flags().StringVarP(&dbType, "type", "t", "mysql", "Specify the type of database; use `mysql` or `psql`. [*Required]")
	dbCreateCmd.Flags().StringVarP(&dbOwner, "owner", "o", "postgres", "Specify owner only for postgres database")
	dbCreateCmd.Flags().StringVarP(&dbUser, "user", "u", "", "Database username [*Required]")
	dbCreateCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbDeleteCmd.Flags().StringVarP(&dbType, "type", "t", "", "Specify the database type for deletion [*Required]")
	dbDeleteCmd.Flags().StringVarP(&dbUser, "user", "u", "", "Database username [*Required]")
	dbDeleteCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbRestoreCmd.Flags().StringVarP(&backupFullFilePath, "backup", "b", "", "Path to the backup file [*Required]")
	dbRestoreCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "Time to schedule the restoration (in a format like HH:MM or a cron-like string)")
	dbRestoreCmd.Flags().StringVarP(&dbType, "type", "t", "", "Specify the database type for restoration [*Required]")
	dbRestoreCmd.Flags().StringVarP(&dbUser, "user", "u", "", "Database username [*Required]")
	dbRestoreCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbListCmd.MarkFlagsRequiredTogether("type", "user")

	dbDeleteCmd.MarkFlagsRequiredTogether("type", "user")

	dbCreateCmd.MarkFlagsRequiredTogether("type", "user")

	dbRestoreCmd.MarkFlagRequired("backup")
	dbRestoreCmd.MarkFlagsRequiredTogether("type", "user")
}
