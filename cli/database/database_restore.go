package database

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Kshitiz-Mhto/xnap/utility"
	"github.com/Kshitiz-Mhto/xnap/utility/common"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbRestoreCmd = &cobra.Command{
	Use:     "restore",
	Aliases: []string{"reset", "restores"},
	Short:   "Restore a database",
	Example: "xnap db restore <database-name> --type <db-type> --user <db_user> --password --backup <path/to/backup-filename> --schedule <schedule-time>",
	Args:    cobra.ExactArgs(1),
	Run:     runRestoreCommand,
}

func runRestoreCommand(cmd *cobra.Command, args []string) {
	databaseName = args[0]
	backupFilePath, _ := cmd.Flags().GetString("backup")
	schedule, _ := cmd.Flags().GetString("schedule")

	if promptPass {
		dbPassword = common.PromptForPassword()
	} else {
		utility.Error("Please include  password paramater `-p`.")
		os.Exit(1)
	}

	switch dbType {
	case "mysql":
		runRestoreCommandForMySQL(backupFilePath, schedule)
	case "psql", "postgres":
		runRestoreCommandForPSQL(backupFilePath, schedule)
	default:
		utility.Error("UnsuppError: Failed to restore database: database is closedorted database type: %s. Use 'mysql', or 'postgres'.\n", dbType)
		os.Exit(1)
	}

}

func runRestoreCommandForMySQL(backupFilePath string, schedule string) {

	// Validate the backup file
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		utility.Error("Backup file does not exist at the specified path: %s", backupFilePath)
		os.Exit(1)
	}

	if schedule == "" {
		performRestoreForMysql(databaseName, backupFilePath)
	} else {
		scheduleRestore(databaseName, backupFilePath, schedule)
	}
}

func runRestoreCommandForPSQL(backupFilePath string, schedule string) {

	// Validate the backup file
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		utility.Error("Backup file does not exist at the specified path: %s", backupFilePath)
		os.Exit(1)
	}

	if schedule == "" {
		performRestoreForPSQL(databaseName, backupFilePath)
	} else {
		scheduleRestore(databaseName, backupFilePath, schedule)
	}
}

func performRestoreForMysql(databaseName, backupFilePath string) {
	var dbExist int = 0
	utility.Info("Starting restore process for database '%s' from file '%s'\n", databaseName, backupFilePath)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}

	if err = db.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", databaseName).Scan(&dbExist).Error; err != nil {
		utility.Error("Failed to check if database exists: %v", err)
		os.Exit(1)
	}

	if dbExist == 0 {
		utility.Info("Database '%s' does not exist. Creating it now...\n\n", utility.Yellow(databaseName))
		if err := db.Exec(fmt.Sprintf("CREATE DATABASE `%s`", databaseName)).Error; err != nil {
			utility.Error("Failed to create database '%s': %v", databaseName, err)
			os.Exit(1)
		}
		utility.Success("Database '%s' created successfully.\n", databaseName)
	}

	command := exec.Command("mysql", "-u", dbUser, "-p"+dbPassword, databaseName)
	backupFile, err := os.Open(backupFilePath)
	if err != nil {
		utility.Error("Failed to open backup file: %v", err)
		os.Exit(1)
	}
	defer backupFile.Close()

	command.Stdin = backupFile

	if err := command.Run(); err != nil {
		utility.Error("Restore process failed: %v", err)
		os.Exit(1)
	}
	utility.Success("Restore process completed successfully.")
}

func performRestoreForPSQL(databaseName, backupFilePath string) {
	utility.Info("Starting restore process for database '%s' from file '%s'\n", databaseName, backupFilePath)

	command := exec.Command("psql", "-u", dbUser, "-p"+dbPassword, databaseName)
	backupFile, err := os.Open(backupFilePath)
	if err != nil {
		utility.Error("Failed to open backup file: %v", err)
		os.Exit(1)
	}
	defer backupFile.Close()

	command.Stdin = backupFile

	if err := command.Run(); err != nil {
		utility.Error("Restore process failed: %v", err)
		os.Exit(1)
	}

	utility.Info("Restore process completed successfully.")
}

func scheduleRestore(databaseName, backupFilePath, schedule string) {
	// Parse the schedule time
	var scheduleTime time.Time
	if len(schedule) == 5 && schedule[2] == ':' { // HH:MM format
		now := time.Now()
		parsedTime, err := time.Parse("15:04", schedule)
		if err != nil {
			utility.Error("Invalid schedule time format. Use 'HH:MM'.")
			os.Exit(1)
		}
		scheduleTime = time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())
		if scheduleTime.Before(now) {
			scheduleTime = scheduleTime.Add(24 * time.Hour) // Schedule for the next day if time has already passed today
		}
	} else {
		utility.Error("Invalid schedule format. Use 'HH:MM'.")
		os.Exit(1)
	}

	utility.Info("Restoration scheduled for %s\n", scheduleTime)

	WG.Add(1)
	// Context to handle cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Wait until the scheduled time
	waitDuration := time.Until(scheduleTime)
	time.AfterFunc(waitDuration, func() {
		defer WG.Done() // Mark the backup task as done
		utility.Info("Starting scheduled restore process at %s\n", time.Now())
		if dbType == "mysql" {
			performRestoreForMysql(databaseName, backupFilePath)
		} else if dbType == "psql" || dbType == "postgres" {
			performRestoreForPSQL(databaseName, backupFilePath)
		}
	})

	// Wait for the backup to finish or for cancellation
	done := make(chan struct{})
	go func() {
		WG.Wait()
		close(done)
	}()

	select {
	case <-done:
		utility.Info("Restore completed successfully.")
	case <-ctx.Done():
		utility.Info("Restore process was canceled.")
	}
}
