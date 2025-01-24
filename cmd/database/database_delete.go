package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/Kshitiz-Mhto/dsync/utility"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm", "del"},
	Short:   "Delete the database",
	Example: "dsync db delete <databse name> --type <database type>",
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
	switch dbType {
	case "mysql":
		deleteMySQLDatabase()
	case "postgres", "psql":
		deletePostgresDatabase()
	default:
		utility.Error("UnsuppError: Failed to create database: sql: database is closedorted database type: %s. Use 'all', 'mysql', or 'postgres'.\n", dbType)
		os.Exit(1)
	}
}

func deleteMySQLDatabase() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/", MySQL_DB_USER, MySQL_DB_PASSWORD, MySQL_DB_HOST, MySQL_DB_PORT)

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

	utility.CloseDBConnection(db)

	fmt.Print("Database deleted successfully !!\n")
}

func deletePostgresDatabase() {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s", POSTGRES_DB_HOSTOST, POSTGRES_DB_PORT, POSTGRES_DB_USER, POSTGRES_DB_PASSWORD)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PostgreSQL: %s\n", err)
		os.Exit(1)
	}

	psql := fmt.Sprintf("DROP DATABASE %s", dbName)
	if err := db.Exec(psql).Error; err != nil {
		utility.Error("Failed to delete database: %v", err)
		os.Exit(1)
	}

	utility.CloseDBConnection(db)
	fmt.Print("Database deleted successfully !!\n")
}
