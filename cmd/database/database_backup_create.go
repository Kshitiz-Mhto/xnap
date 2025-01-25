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
	backupFile.WriteString("/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;\n\n")

	backupFile.WriteString(fmt.Sprintf("/*  --Dump file generated at %s-- */\n", time.Now().Format(time.RFC1123)))

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

	// Write the header of the dump (similar to pg_dump)
	backupFile.WriteString("-- PostgreSQL database dump\n\n")
	backupFile.WriteString(fmt.Sprintf("-- Dumped from database version %s\n", "14.8")) // Adjust version as needed
	backupFile.WriteString("\n")
	backupFile.WriteString("SET statement_timeout = 0;\n")
	backupFile.WriteString("SET lock_timeout = 0;\n")
	backupFile.WriteString("SET idle_in_transaction_session_timeout = 0;\n")
	backupFile.WriteString("SET client_encoding = 'UTF8';\n")
	backupFile.WriteString("SET standard_conforming_strings = on;\n")
	backupFile.WriteString("SELECT pg_catalog.set_config('search_path', '', false);\n")
	backupFile.WriteString("SET check_function_bodies = false;\n")
	backupFile.WriteString("SET xmloption = content;\n")
	backupFile.WriteString("SET client_min_messages = warning;\n")
	backupFile.WriteString("SET row_security = off;\n\n")

	backupFile.WriteString("SET default_tablespace = '';\n")
	backupFile.WriteString("SET default_table_access_method = heap;\n\n")

	// Fetch all tables
	tables := []string{}
	if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname='public'").Scan(&tables).Error; err != nil {
		utility.Error("Error fetching table list: %v", err)
		os.Exit(1)
	}

	// Loop through tables to dump schema, sequences, and data
	for _, table := range tables {
		// Backup schema (CREATE TABLE and related statements) - If not skipping CREATE info
		if !noCreateInfo {
			backupFile.WriteString(fmt.Sprintf("-- Name: %s; Type: TABLE; Schema: public; Owner: postgres\n", table))
			// CREATE TABLE
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

			// Backup sequences (if any)
			var sequenceName string
			seqQuery := fmt.Sprintf(`
				SELECT c.relname 
				FROM pg_class c 
				WHERE c.relkind = 'S' 
				AND c.relname LIKE '%s%%';`, table)
			if err := db.Raw(seqQuery).Scan(&sequenceName).Error; err == nil && sequenceName != "" {
				backupFile.WriteString(fmt.Sprintf("-- Name: %s_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres\n", table))
				backupFile.WriteString(fmt.Sprintf("CREATE SEQUENCE public.%s_id_seq\n", table))
				backupFile.WriteString("    AS integer\n    START WITH 1\n    INCREMENT BY 1\n    NO MINVALUE\n    NO MAXVALUE\n    CACHE 1;\n\n")
			}

			// Backup sequence ownership
			backupFile.WriteString(fmt.Sprintf("-- Name: %s_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres\n", table))
			backupFile.WriteString(fmt.Sprintf("ALTER SEQUENCE public.%s_id_seq OWNED BY public.%s.id;\n\n", table, table))

			// Backup default values for columns
			backupFile.WriteString(fmt.Sprintf("-- Name: %s id; Type: DEFAULT; Schema: public; Owner: postgres\n", table))
			backupFile.WriteString(fmt.Sprintf("ALTER TABLE ONLY public.%s ALTER COLUMN id SET DEFAULT nextval('public.%s_id_seq'::regclass);\n\n", table, table))
		}

		// Backup data using COPY - If not skipping data
		if !noData {
			backupFile.WriteString(fmt.Sprintf("-- Data for Name: %s; Type: TABLE DATA; Schema: public; Owner: postgres\n", table))
			backupFile.WriteString(fmt.Sprintf("COPY public.%s (id, first_name, last_name, role) FROM stdin;\n", table))

			// Fetch the table data and write as COPY format
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
						rowData = append(rowData, fmt.Sprintf("%v", val))
					}
				}
				backupFile.WriteString(strings.Join(rowData, "\t") + "\n")
			}
			backupFile.WriteString("\\.\n\n")
		}

		// Sequence value settings
		if !noCreateInfo {
			backupFile.WriteString(fmt.Sprintf("-- Name: %s_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres\n", table))
			backupFile.WriteString(fmt.Sprintf("SELECT pg_catalog.setval('public.%s_id_seq', 1, true);\n\n", table))
		}

		// Constraints (e.g., PRIMARY KEY) - If not skipping CREATE info
		if !noCreateInfo {
			backupFile.WriteString(fmt.Sprintf("-- Name: %s_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres\n", table))
			backupFile.WriteString(fmt.Sprintf("ALTER TABLE ONLY public.%s ADD CONSTRAINT %s_pkey PRIMARY KEY (id);\n\n", table, table))
		}
	}

	backupFile.WriteString("\n-- PostgreSQL database dump complete\n")
	utility.Success("Backup completed successfully. File saved at: %s", utility.Yellow(backupFullFilePath))
}
