package database

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Kshitiz-Mhto/dsync/utility"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	databaseName       string
	backupFullFilePath string
	noData             bool
	noCreateInfo       bool
)

var dbBackupCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "dysnc db backup create <db-name> --type <db-type> --name <backup-filename> --path <path/to/> --no-data --no-create-info",
	Short:   "Create a new database backup",
	Args:    cobra.ExactArgs(1),
	Run:     dbCreateDatabaseBackup,
}

func dbCreateDatabaseBackup(cmd *cobra.Command, args []string) {
	databaseName = args[0]
	backupFileName, _ = cmd.Flags().GetString("name")
	backupFileNamePath, _ = cmd.Flags().GetString("path")
	noData, _ = cmd.Flags().GetBool("no-data")
	noCreateInfo, _ = cmd.Flags().GetBool("no-create-info")

	if backupFileName == "" {
		// YYYYMMDD_HHMMSS_databaseName_backup.sql
		backupFileName = fmt.Sprintf("%s_%s_backup.sql", time.Now().Format("20060102_150405"), databaseName)
	}

	var err error

	// Resolve the backup file path
	backupFullFilePath, err = filepath.Abs(filepath.Join(backupFileNamePath, backupFileName))
	if err != nil {
		utility.Error("Error resolving backup file path: %v", err)
		os.Exit(1)
	}

	switch dbType {
	case "mysql":
		dbCreateMySQLDatabaseBackup()
	case "postgres", "psql":
		dbCreatePSQLDatabaseBackup()
	default:
		utility.Error("UnsuppError: Failed to create database backup: database is closedorted database type: %s. Use 'mysql', or 'postgres'.\n", dbType)
		os.Exit(1)
	}

}

func dbCreateMySQLDatabaseBackup() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", MySQL_DB_USER, MySQL_DB_PASSWORD, MySQL_DB_HOST, MySQL_DB_PORT, databaseName)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	utility.Info("Starting backup for %s database '%s'...", utility.Yellow(dbType), utility.Yellow(databaseName))

	// Open the backup file
	backupFile, err := os.Create(backupFullFilePath)
	if err != nil {
		utility.Error("Error creating backup file: %v", err)
		os.Exit(1)
	}
	defer backupFile.Close()

	// Write header information
	backupFile.WriteString(fmt.Sprintf("-- dSync-MySQL dump v1.0.0, for %s/%s\n", runtime.GOOS, runtime.GOARCH))
	backupFile.WriteString(fmt.Sprintf("--\n-- Host: %s    Database: %s\n--\n", MySQL_DB_HOST, databaseName))
	backupFile.WriteString("-- ------------------------------------------------------\n")

	backupFile.WriteString("/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;\n")
	backupFile.WriteString("/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;\n")
	backupFile.WriteString("/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;\n")
	backupFile.WriteString("/*!50503 SET NAMES utf8mb4 */;\n")
	backupFile.WriteString("/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;\n")
	backupFile.WriteString("/*!40103 SET TIME_ZONE='+00:00' */;\n")
	backupFile.WriteString("/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;\n")
	backupFile.WriteString("/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;\n")
	backupFile.WriteString("/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;\n")
	backupFile.WriteString("/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;\n\n")

	// Fetch all table names
	tables := []string{}
	if err := db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		utility.Error("Error fetching table list: %v", err)
		os.Exit(1)
	}

	for _, table := range tables {
		// Lock the table before any operation
		if err := db.Exec(fmt.Sprintf("LOCK TABLES `%s` READ", table)).Error; err != nil {
			utility.Error("Error locking table '%s': %v", table, err)
			os.Exit(1)
		}

		if !noCreateInfo {
			// Backup table schema
			backupFile.WriteString(fmt.Sprintf("--\n-- Table structure for table `%s`\n--\n\n", table))
			backupFile.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", table))
			backupFile.WriteString("/*!40101 SET @saved_cs_client     = @@character_set_client */;\n")
			backupFile.WriteString("/*!50503 SET character_set_client = utf8mb4 */;\n")
			var createTableQuery string
			if err := db.Raw(fmt.Sprintf("SHOW CREATE TABLE `%s`", table)).Row().Scan(&table, &createTableQuery); err != nil {
				utility.Error("Error fetching schema for table '%s': %v", table, err)
				os.Exit(1)
			}
			backupFile.WriteString(createTableQuery + ";\n\n")
		}

		if !noData {
			// Backup table data
			backupFile.WriteString(fmt.Sprintf("--\n-- Dumping data for table `%s`\n--\n\n", table))
			backupFile.WriteString(fmt.Sprintf("LOCK TABLES `%s` WRITE;\n", table))
			backupFile.WriteString(fmt.Sprintf("/*!40000 ALTER TABLE `%s` DISABLE KEYS */;\n", table))
			rows, err := db.Raw(fmt.Sprintf("SELECT * FROM `%s`", table)).Rows()
			if err != nil {
				utility.Error("Error fetching data for table '%s': %v", table, err)
				os.Exit(1)
			}
			columns, _ := rows.Columns()
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			for rows.Next() {
				rows.Scan(valuePtrs...)
				rowData := []string{}
				for _, val := range values {
					if val == nil {
						rowData = append(rowData, "NULL")
					} else {
						switch v := val.(type) {
						case []byte: // Convert byte slices to string
							rowData = append(rowData, fmt.Sprintf("'%s'", string(v)))
						case string:
							rowData = append(rowData, fmt.Sprintf("'%s'", v))
						case int, int64, float64: // Handle numeric types
							rowData = append(rowData, fmt.Sprintf("%v", v))
						default: // Handle other types generically
							rowData = append(rowData, fmt.Sprintf("'%v'", v))
						}
					}
				}
				insertQuery := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s);\n", table, strings.Join(columns, ", "), strings.Join(rowData, ", "))
				backupFile.WriteString(insertQuery)
			}
			backupFile.WriteString(fmt.Sprintf("/*!40000 ALTER TABLE `%s` ENABLE KEYS */;\n", table))
			backupFile.WriteString("UNLOCK TABLES;\n\n")
			backupFile.WriteString("\n")
		}

		// Unlock the table after operation
		if err := db.Exec("UNLOCK TABLES").Error; err != nil {
			utility.Error("Error unlocking tables: %v", err)
			os.Exit(1)
		}
	}

	// Footer
	backupFile.WriteString("/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;\n")
	backupFile.WriteString("/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;\n")
	backupFile.WriteString("/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;\n")
	backupFile.WriteString("/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;\n")
	backupFile.WriteString("/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;\n")
	backupFile.WriteString("/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;\n")
	backupFile.WriteString("/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;\n")
	backupFile.WriteString("/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;\n")

	utility.Success("Backup completed successfully. File saved at: %s", utility.Yellow(backupFullFilePath))
}

func dbCreatePSQLDatabaseBackup() {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", POSTGRES_DB_HOSTOST, POSTGRES_DB_PORT, POSTGRES_DB_USER, POSTGRES_DB_PASSWORD, databaseName)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PostgreSQL: %s", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)

	utility.Info("Starting backup for %s database '%s'...", utility.Yellow(dbType), utility.Yellow(databaseName))

	// Open the backup file
	backupFile, err := os.Create(backupFullFilePath)
	if err != nil {
		utility.Error("Error creating backup file: %v", err)
		os.Exit(1)
	}
	defer backupFile.Close()

	tables := []string{}
	if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname='public'").Scan(&tables).Error; err != nil {
		utility.Error("Error fetching table list: %v", err)
		os.Exit(1)
	}

	backupFile.WriteString(fmt.Sprintf("-- Host: %s/%s\n", POSTGRES_DB_HOSTOST, POSTGRES_DB_PORT))
	backupFile.WriteString(fmt.Sprintf("-- Database: %s\n\n", databaseName))

	for _, table := range tables {
		if !noCreateInfo {
			// Backup table schema using pg_catalog
			backupFile.WriteString(fmt.Sprintf("-- Schema for table '%s'\n", table))

			var createTableQuery string
			query := fmt.Sprintf(`
				SELECT 'CREATE TABLE ' || relname || E'\n(\n' ||
				array_to_string(
					array_agg(
						'    ' || column_name || ' ' || type_name || 
						CASE WHEN is_nullable THEN ' NULL' ELSE ' NOT NULL' END
					), E',\n'
				) || E'\n);\n'
				FROM (
					SELECT
						c.relname,
						a.attname AS column_name,
						pg_catalog.format_type(a.atttypid, a.atttypmod) AS type_name,
						a.attnotnull = false AS is_nullable
					FROM
						pg_class c
					JOIN
						pg_attribute a ON a.attrelid = c.oid
					WHERE
						c.relname = '%s'
						AND a.attnum > 0 -- This filters out system columns
				) sub
				GROUP BY relname;
			`, table)

			if err := db.Raw(query).Scan(&createTableQuery).Error; err != nil {
				utility.Error("Error fetching schema for table '%s': %v", table, err)
				os.Exit(1)
			}
			backupFile.WriteString(createTableQuery + ";\n\n")
		}

		if !noData {
			// Backup table data
			backupFile.WriteString(fmt.Sprintf("-- Data for table '%s'\n", table))
			rows, err := db.Raw(fmt.Sprintf("SELECT * FROM \"%s\"", table)).Rows()
			if err != nil {
				utility.Error("Error fetching data for table '%s': %v", table, err)
				os.Exit(1)
			}

			columns, _ := rows.Columns()
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			for rows.Next() {
				rows.Scan(valuePtrs...)
				rowData := []string{}
				for _, val := range values {
					if val == nil {
						rowData = append(rowData, "NULL")
					} else {
						rowData = append(rowData, fmt.Sprintf("'%v'", val))
					}
				}
				insertQuery := fmt.Sprintf("INSERT INTO \"%s\" (%s) VALUES (%s);\n", table, strings.Join(columns, ", "), strings.Join(rowData, ", "))
				backupFile.WriteString(insertQuery)
			}
			backupFile.WriteString("\n")
		}
	}
	utility.Success("Backup completed successfully. File saved at: %s", utility.Yellow(backupFullFilePath))

}
