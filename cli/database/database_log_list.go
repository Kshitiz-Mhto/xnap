package database

import (
	"fmt"
	"os"

	"github.com/Kshitiz-Mhto/xnap/cli/logs"
	"github.com/Kshitiz-Mhto/xnap/pkg/config"
	"github.com/Kshitiz-Mhto/xnap/utility"
	"github.com/Kshitiz-Mhto/xnap/utility/common"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbLogsListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "list all the log list",
	Example: "xnap db log ls --type <type> --user <db_user> --password",
	Run:     runDbLogsList,
}

func runDbLogsList(cmd *cobra.Command, args []string) {
	if promptPass {
		dbPassword = common.PromptForPassword()
	} else {
		utility.Error("Please include  password paramater `-p`.")
		os.Exit(1)
	}

	// Switch between database types
	switch dbType {
	case "mysql":
		dbLogsListFromMySQL()
	case "postgres", "psql":
		dbLogsListFromPSQL()
	default:
		utility.Error("Unsupported database type: %s. Use 'mysql', or 'postgres'.", dbType)
		os.Exit(1)
	}
}

func dbLogsListFromMySQL() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, config.Envs.XNAP_DB)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}

	var listOfLogs []logs.Log

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at ASC", config.Envs.XNAP_LOGS_TABLE)
	if err = db.Raw(query).Scan(&listOfLogs).Error; err != nil {
		utility.Error("Failed to fetch rows from MySQL: %s", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	utility.Info("Listing the logs list from table `%s`", config.Envs.XNAP_LOGS_TABLE)

	ow := utility.NewOutputWriter()
	oi := utility.NewOutputWriter()

	oi.AppendDataWithLabel("db_host", MySQL_DB_HOST, "DB_HOST")
	oi.AppendDataWithLabel("port", MySQL_DB_PORT, "DB_PORT")
	oi.AppendDataWithLabel("user", dbUser, "DB_USER")
	oi.AppendDataWithLabel("type", "MySQL", "DB_TYPE")
	oi.FinishAndPrintOutput()

	if listOfLogs == nil {
		utility.Warning("Failed to fetch populated list of backups or empty rows")
		os.Exit(1)
	}
	if len(listOfLogs) == 0 {
		utility.Error("No backups found in the database.")
		os.Exit(1)
	}

	for _, backup := range listOfLogs {
		ow.StartLine()

		ow.AppendDataWithLabel("id", backup.ID.String(), "id")
		ow.AppendDataWithLabel("action", backup.Action, "action")
		ow.AppendDataWithLabel("command", backup.Command, "command")
		ow.AppendDataWithLabel("status", backup.Status, "status")
		ow.AppendDataWithLabel("error_message", backup.ErrorMessage, "error_message")
		ow.AppendDataWithLabel("user_name", backup.UserName, "user_name")
		ow.AppendDataWithLabel("execution_duration", fmt.Sprintf("%.2f", backup.ExecutionDuration), "execution_duration")
		ow.AppendDataWithLabel("created_at", backup.CreatedAt.String(), "created_at")
		ow.AppendDataWithLabel("updated_at", backup.UpdatedAt.String(), "updated_at")
	}
	ow.FinishAndPrintOutput()
}

func dbLogsListFromPSQL() {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", POSTGRES_DB_HOST, POSTGRES_DB_PORT, dbUser, dbPassword, config.Envs.XNAP_DB)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PostgreSQL: %s", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	var listOfLogs []logs.Log

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at ASC", config.Envs.XNAP_LOGS_TABLE)
	if err = db.Raw(query).Scan(&listOfLogs).Error; err != nil {
		utility.Error("Failed to fetch rows from MySQL: %s", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	utility.Info("Listing the logs list from table `%s`", config.Envs.XNAP_LOGS_TABLE)

	ow := utility.NewOutputWriter()
	oi := utility.NewOutputWriter()

	oi.AppendDataWithLabel("db_host", POSTGRES_DB_HOST, "DB_HOST")
	oi.AppendDataWithLabel("port", POSTGRES_DB_PORT, "DB_PORT")
	oi.AppendDataWithLabel("user", dbUser, "DB_USER")
	oi.AppendDataWithLabel("type", "Postgresql", "DB_TYPE")
	oi.FinishAndPrintOutput()

	if listOfLogs == nil {
		utility.Warning("Failed to fetch populated list of backups or empty rows")
		os.Exit(1)
	}
	if len(listOfLogs) == 0 {
		utility.Error("No backups found in the database.")
		os.Exit(1)
	}

	for _, backup := range listOfLogs {
		ow.StartLine()

		ow.AppendDataWithLabel("id", backup.ID.String(), "id")
		ow.AppendDataWithLabel("action", backup.Action, "action")
		ow.AppendDataWithLabel("command", backup.Command, "command")
		ow.AppendDataWithLabel("status", backup.Status, "status")
		ow.AppendDataWithLabel("error_message", backup.ErrorMessage, "error_message")
		ow.AppendDataWithLabel("user_name", backup.UserName, "user_name")
		ow.AppendDataWithLabel("execution_duration", fmt.Sprintf("%.2f", backup.ExecutionDuration), "execution_duration")
		ow.AppendDataWithLabel("created_at", backup.CreatedAt.String(), "created_at")
		ow.AppendDataWithLabel("updated_at", backup.UpdatedAt.String(), "updated_at")
	}
	ow.FinishAndPrintOutput()

}
