package database

import (
	"fmt"
	"os"

	"github.com/Kshitiz-Mhto/dsync/pkg/config"
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
	Example: "dsync db ls",
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
	case "postgres":
		listPostgresDatabases()
	default:
		utility.Error("Unsupported database type: %s. Use 'all', 'mysql', or 'postgres'.", dbType)
		os.Exit(1)
	}

}

func listMySQLDatabases() {

	var (
		DB_USER     string = config.Envs.MySQL_DB_USER
		DB_PASSWORD string = config.Envs.MySQL_DB_PASSWORD
		HOST        string = config.Envs.MySQL_DB_HOST
		PORT        string = config.Envs.MySQL_DB_PORT
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", DB_USER, DB_PASSWORD, HOST, PORT)

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

	ow := utility.NewOutputWriter()
	oi := utility.NewOutputWriter()

	oi.AppendDataWithLabel("db_host", HOST, "DB_HOST")
	oi.AppendDataWithLabel("port", PORT, "DB_PORT")
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

	var (
		DB_USER     string = config.Envs.POSTGRES_DB_USER
		DB_PASSWORD string = config.Envs.POSTGRES_DB_PASSWORD
		HOST        string = config.Envs.POSTGRES_DB_HOST
		PORT        string = config.Envs.POSTGRES_DB_PORT
	)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", HOST, PORT, DB_USER, DB_PASSWORD)

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

	ow := utility.NewOutputWriter()
	oi := utility.NewOutputWriter()

	oi.AppendDataWithLabel("db_host", HOST, "DB_HOST")
	oi.AppendDataWithLabel("port", PORT, "DB_PORT")
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
