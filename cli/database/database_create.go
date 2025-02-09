package database

import (
	"fmt"
	"os"

	"github.com/Kshitiz-Mhto/xnap/utility"
	"github.com/Kshitiz-Mhto/xnap/utility/common"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "xnap db create <db name> --type <db_type> --user <db_user> --password",
	Short:   "Create new database",
	Args:    cobra.ExactArgs(1),
	Run:     dbCreation,
}

func dbCreation(cmd *cobra.Command, args []string) {
	dbName = args[0]

	if promptPass {
		dbPassword = common.PromptForPassword()
	} else {
		utility.Error("Please include  password paramater `-p`.")
		os.Exit(1)
	}

	switch dbType {
	case "mysql":
		createMySQLDatabase()
	case "postgres", "psql":
		createPostgresDatabase()
	default:
		utility.Error("UnsuppError: Failed to create database: sql: database is closedorted database type: %s. Use 'all', 'mysql', or 'postgres'.\n", dbType)
		os.Exit(1)
	}
}

func createMySQLDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}
	sql := fmt.Sprintf("CREATE DATABASE %s CHARACTER SET %s COLLATE %s;", dbName, "utf8mb4", "utf8mb4_general_ci")
	if err := db.Exec(sql).Error; err != nil {
		utility.Error("Failed to create database: %v", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	ow := utility.NewOutputWriter()

	ow.StartLine()
	ow.AppendDataWithLabel("mysql_db_name", dbName, "DB_Name")
	ow.FinishAndPrintOutput()
	utility.Success("Database created successfully !!")
}

func createPostgresDatabase() {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s", POSTGRES_DB_HOST, POSTGRES_DB_PORT, dbUser, dbPassword)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PostgreSQL: %s", err)
		os.Exit(1)
	}

	psql := fmt.Sprintf("CREATE DATABASE %s WITH OWNER = %s ENCODING = '%s' LC_COLLATE = '%s' LC_CTYPE = '%s' CONNECTION LIMIT = -1;", dbName, dbOwner, "UTF-8", "en_US.UTF-8", "en_US.UTF-8")
	if err := db.Exec(psql).Error; err != nil {
		utility.Error("Failed to create database: %v", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	ow := utility.NewOutputWriter()

	ow.StartLine()
	ow.AppendDataWithLabel("pssql_db_name", dbName, "DB_Name")
	ow.FinishAndPrintOutput()
	utility.Success("Database created successfully !!")
}
