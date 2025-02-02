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

var localBackupListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "list all the backup list",
	Example: "xnap local backup ls --type <type> --user <db_user> --password",
	Run:     runLocalBackupList,
}

func runLocalBackupList(cmd *cobra.Command, args []string) {
	if promptPass {
		dbPassword = common.PromptForPassword()
	} else {
		utility.Error("Please include  password paramater `-p`.")
		os.Exit(1)
	}

	// Switch between database types
	switch dbType {
	case "mysql":
		listBackupListFromMySQL()
	case "postgres", "psql":
		listBackupListFromPSQL()
	default:
		utility.Error("Unsupported database type: %s. Use 'mysql', or 'postgres'.", dbType)
		os.Exit(1)
	}
}

func listBackupListFromMySQL() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, config.Envs.XNAP_DB)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}

	var listOfBackups []Backup

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at ASC", config.Envs.XNAP_BACKUP_TABLE)
	if err = db.Raw(query).Scan(&listOfBackups).Error; err != nil {
		utility.Error("Failed to fetch rows from MySQL: %s", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	utility.Info("Listing the backup list from table `%s`", config.Envs.XNAP_BACKUP_TABLE)

	ow := utility.NewOutputWriter()
	oi := utility.NewOutputWriter()

	oi.AppendDataWithLabel("db_host", MySQL_DB_HOST, "DB_HOST")
	oi.AppendDataWithLabel("port", MySQL_DB_PORT, "DB_PORT")
	oi.AppendDataWithLabel("user", dbUser, "DB_USER")
	oi.AppendDataWithLabel("type", "MySQL", "DB_TYPE")
	oi.FinishAndPrintOutput()

	if listOfBackups == nil {
		utility.Warning("Failed to fetch populated list of backups or empty rows")
		os.Exit(1)
	}
	if len(listOfBackups) == 0 {
		utility.Error("No backups found in the database.")
		os.Exit(1)
	}

	for _, backup := range listOfBackups {
		ow.StartLine()

		ow.AppendDataWithLabel("id", backup.ID.String(), "id")
		ow.AppendDataWithLabel("file_name", backup.FileName, "file_name")
		ow.AppendDataWithLabel("source_path", backup.SourcePath, "source_path")
		ow.AppendDataWithLabel("backup_path", backup.BackupPath, "backup_path")
		ow.AppendDataWithLabel("og_file_name", backup.OgFileName, "og_file_name")
		ow.AppendDataWithLabel("created_at", backup.CreatedAt.String(), "created_at")
		ow.AppendDataWithLabel("updated_at", backup.UpdatedAt.String(), "updated_at")

	}

	ow.FinishAndPrintOutput()
}

func listBackupListFromPSQL() {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", POSTGRES_DB_HOST, POSTGRES_DB_PORT, dbUser, dbPassword, config.Envs.XNAP_DB)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PostgreSQL: %s", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	var listOfBackups []Backup

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at ASC", config.Envs.XNAP_BACKUP_TABLE)
	if err = db.Raw(query).Scan(&listOfBackups).Error; err != nil {
		utility.Error("Failed to fetch rows from MySQL: %s", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	utility.Info("Listing the backup list from table `%s`", config.Envs.XNAP_BACKUP_TABLE)

	ow := utility.NewOutputWriter()
	oi := utility.NewOutputWriter()

	oi.AppendDataWithLabel("db_host", POSTGRES_DB_HOST, "DB_HOST")
	oi.AppendDataWithLabel("port", POSTGRES_DB_PORT, "DB_PORT")
	oi.AppendDataWithLabel("user", dbUser, "DB_USER")
	oi.AppendDataWithLabel("type", "Postgresql", "DB_TYPE")
	oi.FinishAndPrintOutput()

	if listOfBackups == nil {
		utility.Warning("Failed to fetch populated list of backups or empty rows")
		os.Exit(1)
	}
	if len(listOfBackups) == 0 {
		utility.Error("No backups found in the database.")
		os.Exit(1)
	}

	for _, backup := range listOfBackups {
		ow.StartLine()

		ow.AppendDataWithLabel("id", backup.ID.String(), "id")
		ow.AppendDataWithLabel("file_name", backup.FileName, "file_name")
		ow.AppendDataWithLabel("source_path", backup.SourcePath, "source_path")
		ow.AppendDataWithLabel("backup_path", backup.BackupPath, "backup_path")
		ow.AppendDataWithLabel("og_file_name", backup.OgFileName, "og_file_name")
		ow.AppendDataWithLabel("created_at", backup.CreatedAt.String(), "created_at")
		ow.AppendDataWithLabel("updated_at", backup.UpdatedAt.String(), "updated_at")

	}

	ow.FinishAndPrintOutput()
}
