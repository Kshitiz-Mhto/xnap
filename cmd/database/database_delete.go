package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/Kshitiz-Mhto/xnap/utility"
	"github.com/Kshitiz-Mhto/xnap/utility/common"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm", "del"},
	Short:   "Delete the database",
	Example: "xnap db delete <databse name> --type <database_type> --user <db_user> --password",
	Args:    cobra.ExactArgs(1),
	Run:     dbDeletion,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if dbType == "" {
			return errors.New("-t (type) flag is required for deletion of database")
		}
		return nil
	},
}

func dbDeletion(cmd *cobra.Command, args []string) {
	dbName = args[0]

	if promptPass {
		dbPassword = common.PromptForPassword()
	} else {
		utility.Error("Please include  password paramater `-p`.")
		os.Exit(1)
	}

	switch dbType {
	case "mysql":
		deleteMySQLDatabase()
	case "postgres", "psql":
		deletePostgresDatabase()
	default:
		utility.Error("UnsuppError: Failed to delete database: database is closedorted database type: %s. Use 'mysql', or 'postgres'.", dbType)
		os.Exit(1)
	}
}

func deleteMySQLDatabase() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}
	sql := fmt.Sprintf("DROP DATABASE %s", dbName)
	if err := db.Exec(sql).Error; err != nil {
		utility.Error("Failed to delete database: %v", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	utility.Success("Database deleted successfully !!")
}

func deletePostgresDatabase() {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s", POSTGRES_DB_HOST, POSTGRES_DB_PORT, dbUser, dbPassword)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PostgreSQL: %s", err)
		os.Exit(1)
	}

	psql := fmt.Sprintf("DROP DATABASE %s", dbName)
	if err := db.Exec(psql).Error; err != nil {
		utility.Error("Failed to delete database: %v", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)
	utility.Success("Database deleted successfully !!")
}
