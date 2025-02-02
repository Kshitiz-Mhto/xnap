package local

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

var localBackupListDeletionCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"del", "delete", "rm"},
	Short:   "Remove all the rows from backup table",
	Example: "xnap local backup rm --type <type> --user <db_user> --password",
	Run:     runLocalBackupListDeletion,
}

func runLocalBackupListDeletion(cmd *cobra.Command, args []string) {
	if promptPass {
		dbPassword = common.PromptForPassword()
	} else {
		utility.Error("Please include  password paramater `-p`.")
		os.Exit(1)
	}

	// Switch between database types
	switch dbType {
	case "mysql":
		deleteBackupListFromMySQL()
	case "postgres", "psql":
		deleteBackupListFromPSQL()
	default:
		utility.Error("Unsupported database type: %s. Use 'mysql', or 'postgres'.", dbType)
		os.Exit(1)
	}
}
func deleteBackupListFromMySQL() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, config.Envs.XNAP_DB)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}
	defer utility.CloseDBConnection(db)

	sql := fmt.Sprintf("DELETE FROM %s", config.Envs.XNAP_BACKUP_TABLE)
	if err = db.Exec(sql).Error; err != nil {
		utility.Error("Failed to delete table rows: %v", err)
		os.Exit(1)
	}

	utility.Success("Table rows deleted successfully !!")
}

func deleteBackupListFromPSQL() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", POSTGRES_DB_HOST, POSTGRES_DB_PORT, dbUser, dbPassword, config.Envs.XNAP_DB)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PostgreSQL: %s", err)
		os.Exit(1)
	}
	defer utility.CloseDBConnection(db)

	psql := fmt.Sprintf("DELETE FROM %s", config.Envs.XNAP_BACKUP_TABLE)
	if err = db.Exec(psql).Error; err != nil {
		utility.Error("Failed to delete table rows: %v", err)
		os.Exit(1)
	}

	utility.Success("Table rows deleted successfully !!")

}
