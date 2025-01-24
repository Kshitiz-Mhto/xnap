package database

import (
	"fmt"
	"os"

	"github.com/Kshitiz-Mhto/dsync/utility"
	"github.com/Kshitiz-Mhto/dsync/utility/common"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbType string

var dbListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "list all the databases",
	Example: "dsync db ls --type <type>",
	Run:     listDatabases,
}

func listDatabases(cmd *cobra.Command, args []string) {

	// Switch between database types
	switch dbType {
	case "all":
		listMySQLDatabases()
		listPostgresDatabases()
	case "mysql":
		listMySQLDatabases()
	case "postgres", "psql":
		listPostgresDatabases()
	default:
		utility.Error("Unsupported database type: %s. Use 'all', 'mysql', or 'postgres'.", dbType)
		os.Exit(1)
	}

}

func listMySQLDatabases() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", MySQL_DB_USER, MySQL_DB_PASSWORD, MySQL_DB_HOST, MySQL_DB_PORT)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}
	var databases []string

	if err := db.Raw("SHOW DATABASES").Scan(&databases).Error; err != nil {
		utility.Error("Failed to fetch databases: %v", err)
		os.Exit(1)
	}
	utility.CloseDBConnection(db)

	ow := utility.NewOutputWriter()
	oi := utility.NewOutputWriter()

	oi.AppendDataWithLabel("db_host", MySQL_DB_HOST, "DB_HOST")
	oi.AppendDataWithLabel("port", MySQL_DB_PORT, "DB_PORT")
	oi.AppendDataWithLabel("type", "MySQL", "DB_TYPE")
	oi.FinishAndPrintOutput()

	for _, dbName := range databases {

		ow.StartLine()
		ow.AppendDataWithLabel("mysql_db_name", dbName, "DB_Name")

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("dn_name", dbName, "DB_Name")
		}
	}
	ow.FinishAndPrintOutput()

}

func listPostgresDatabases() {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", POSTGRES_DB_HOSTOST, POSTGRES_DB_PORT, POSTGRES_DB_USER, POSTGRES_DB_PASSWORD)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PostgreSQL: %s", err)
		os.Exit(1)
	}

	var databases []string

	// PostgreSQL uses a different query to list databases
	if err := db.Raw("SELECT datname FROM pg_database WHERE datistemplate = false").Scan(&databases).Error; err != nil {
		utility.Error("Failed to fetch databases: %v", err)
		os.Exit(1)
	}
	utility.CloseDBConnection(db)

	ow := utility.NewOutputWriter()
	oi := utility.NewOutputWriter()

	oi.AppendDataWithLabel("db_host", POSTGRES_DB_HOSTOST, "DB_HOST")
	oi.AppendDataWithLabel("port", POSTGRES_DB_PORT, "DB_PORT")
	oi.AppendDataWithLabel("type", "PostgreSQL", "DB_TYPE")
	oi.FinishAndPrintOutput()

	for _, dbName := range databases {

		ow.StartLine()
		ow.AppendDataWithLabel("postgres_db_name", dbName, "DB_Name")

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("db_name", dbName, "DB_Name")
		}
	}
	ow.FinishAndPrintOutput()
}
