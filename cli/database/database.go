package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/Kshitiz-Mhto/xnap/cli/alert"
	"github.com/Kshitiz-Mhto/xnap/cli/logs"
	"github.com/Kshitiz-Mhto/xnap/pkg/config"
	"github.com/Kshitiz-Mhto/xnap/utility"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	dbListCmd.Flags().StringVarP(&dbType, "type", "t", "", "Database type: mysql or postgres")
	dbListCmd.Flags().StringVarP(&dbUser, "user", "u", dbUser, "Database username")
	dbListCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbCreateCmd.Flags().StringVarP(&dbType, "type", "t", "mysql", "Create database type MySQL")
	dbCreateCmd.Flags().StringVarP(&dbOwner, "owner", "o", "postgres", "Specify owner only for postgres database")
	dbCreateCmd.Flags().StringVarP(&dbUser, "user", "u", dbUser, "Database username")
	dbCreateCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbDeleteCmd.Flags().StringVarP(&dbType, "type", "t", "", "Specify the database type for deletion")
	dbDeleteCmd.Flags().StringVarP(&dbUser, "user", "u", dbUser, "Database username")
	dbDeleteCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbRestoreCmd.Flags().StringVarP(&backupFullFilePath, "backup", "b", "", "Path to the backup file (required)")
	dbRestoreCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "Time to schedule the restoration (in a format like HH:MM or a cron-like string)")
	dbRestoreCmd.Flags().StringVarP(&dbType, "type", "t", "", "Specify the database type for restoration")
	dbRestoreCmd.Flags().StringVarP(&dbUser, "user", "u", dbUser, "Database username")
	dbRestoreCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")

	dbListCmd.MarkFlagsRequiredTogether("type", "user")

	dbDeleteCmd.MarkFlagsRequiredTogether("type", "user")

	dbCreateCmd.MarkFlagsRequiredTogether("type", "user")

	dbRestoreCmd.MarkFlagRequired("backup")
	dbRestoreCmd.MarkFlagsRequiredTogether("type", "user")
}

func logCommand(dbType, dbUser, dbPassword, host, port, action, command, status, errorMessage string, userName string, duration float64) error {
	var err error

	logEntry := &logs.Log{
		Action:            action,
		Command:           command,
		Status:            status,
		ErrorMessage:      errorMessage,
		UserName:          userName,
		ExecutionDuration: duration,
	}

	switch dbType {
	case "mysql":
		err = AddLogToMysql(dbUser, dbPassword, host, port, logEntry)
	case "psql", "postgres":
		err = AddLogtoPSQL(dbUser, dbPassword, host, port, logEntry)
	}

	return err
}

func AddLogToMysql(dbUser, dbPassword, host, port string, logEntry *logs.Log) error {
	var lastLogEntry logs.Log

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, host, port, config.Envs.LOG_DB)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		return err
	}

	if err := db.Create(logEntry).Error; err != nil {
		return err
	}

	if err := db.Order("created_at DESC").First(&lastLogEntry).Error; err != nil {
		fmt.Println("Error fetching the last row:", err)
	}

	if logEntry.Status == config.Envs.BACKUP_OR_RESTORE_STATUS {
		vars := map[string]interface{}{
			"ID":                lastLogEntry.ID,
			"Action":            lastLogEntry.Action,
			"Command":           lastLogEntry.Command,
			"Status":            lastLogEntry.Status,
			"ErrorMessage":      lastLogEntry.ErrorMessage,
			"UserName":          lastLogEntry.UserName,
			"ExecutionDuration": lastLogEntry.ExecutionDuration,
			"CreatedAt":         lastLogEntry.CreatedAt,
			"UpdatedAt":         lastLogEntry.UpdatedAt,
			"dbType":            dbType,
		}
		alert.HTMLTemplateEmailHandler(config.Envs.OWNER_EMAIL, vars)
	}
	utility.Success("Log is enteried successfully!!")

	return nil
}

func AddLogtoPSQL(dbUser, dbPassword, host, port string, logEntry *logs.Log) error {
	var lastLogEntry logs.Log

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, dbUser, dbPassword, config.Envs.LOG_DB)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PSQL: %s", err)
		return err
	}

	if err := db.Create(logEntry).Error; err != nil {
		return err
	}

	if err := db.Order("created_at DESC").First(&lastLogEntry).Error; err != nil {
		return err
	}

	if logEntry.Status == config.Envs.BACKUP_OR_RESTORE_STATUS {
		vars := map[string]interface{}{
			"ID":                lastLogEntry.ID,
			"Action":            lastLogEntry.Action,
			"Command":           lastLogEntry.Command,
			"Status":            lastLogEntry.Status,
			"ErrorMessage":      lastLogEntry.ErrorMessage,
			"UserName":          lastLogEntry.UserName,
			"ExecutionDuration": lastLogEntry.ExecutionDuration,
			"CreatedAt":         lastLogEntry.CreatedAt,
			"UpdatedAt":         lastLogEntry.UpdatedAt,
			"dbType":            dbType,
		}

		alert.HTMLTemplateEmailHandler(config.Envs.OWNER_EMAIL, vars)
	}

	utility.Success("Log is enteried successfully!!")

	return nil
}

func SetFailureStatus(msg string) {
	status = "failure"
	errorMessage = msg
}
