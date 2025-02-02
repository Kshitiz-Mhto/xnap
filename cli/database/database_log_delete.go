package database

import (
	"fmt"
	"os"

	"github.com/Kshitiz-Mhto/xnap/pkg/config"
	"github.com/Kshitiz-Mhto/xnap/utility"
	"github.com/Kshitiz-Mhto/xnap/utility/common"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbLogListDeletionCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"del", "delete", "rm"},
	Short:   "Remove all the rows from logs table",
	Example: "xnap db log rm --type <type> --user <db_user> --password",
	Run:     runDbLogspListDeletion,
}

func runDbLogspListDeletion(cmd *cobra.Command, args []string) {
	if promptPass {
		dbPassword = common.PromptForPassword()
	} else {
		utility.Error("Please include  password paramater `-p`.")
		os.Exit(1)
	}

	// Switch between database types
	switch dbType {
	case "mysql":
		deleteLogsListFromMySQL()
	case "postgres", "psql":
		deleteLogsListFromPSQL()
	default:
		utility.Error("Unsupported database type: %s. Use 'mysql', or 'postgres'.", dbType)
		os.Exit(1)
	}
}
func deleteLogsListFromMySQL() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, config.Envs.XNAP_DB)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}
	defer utility.CloseDBConnection(db)

	utility.Warning("Starting Deletion process for `%s` table", config.Envs.XNAP_LOGS_TABLE)

	sql := fmt.Sprintf("DELETE FROM %s", config.Envs.XNAP_LOGS_TABLE)
	if err = db.Exec(sql).Error; err != nil {
		utility.Error("Failed to delete table rows: %v", err)
		os.Exit(1)
	}

	utility.Success("Table rows deleted successfully !!")
}

func deleteLogsListFromPSQL() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", POSTGRES_DB_HOST, POSTGRES_DB_PORT, dbUser, dbPassword, config.Envs.XNAP_DB)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PostgreSQL: %s", err)
		os.Exit(1)
	}
	defer utility.CloseDBConnection(db)

	utility.Warning("Starting Deletion process for `%s` table", config.Envs.XNAP_LOGS_TABLE)

	psql := fmt.Sprintf("DELETE FROM %s", config.Envs.XNAP_BACKUP_TABLE)
	if err = db.Exec(psql).Error; err != nil {
		utility.Error("Failed to delete table rows: %v", err)
		os.Exit(1)
	}

	utility.Success("Table rows deleted successfully !!")

}
