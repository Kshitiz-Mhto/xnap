package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kshitiz-Mhto/xnap/cli/database"
	"github.com/Kshitiz-Mhto/xnap/pkg/config"
	"github.com/Kshitiz-Mhto/xnap/utility"
	"github.com/Kshitiz-Mhto/xnap/utility/common"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var LocalBackupCmd = &cobra.Command{
	Use:     "backup",
	Short:   "Backup the file from local storage",
	Long:    "Backup the file by copying it to backup locationa and store the file path in database as reference",
	Example: "xnap local backup --type <database_type> -u <username> -p --source </to/path/file> --path </to/path/backup_location> --version <version_number> --schedule <schedule_HH:MM>",
	Run:     runLocalBackupCommand,
}

func runLocalBackupCommand(cmd *cobra.Command, args []string) {
	sourcePath, _ = cmd.Flags().GetString("source")
	backupDirPath, _ = cmd.Flags().GetString("path")
	versionNum, _ = cmd.Flags().GetString("version")
	command = strings.Join(os.Args, " ")
	start = time.Now()
	status = "success"
	errorMessage = ""

	if promptPass {
		dbPassword = common.PromptForPassword()
	} else {
		utility.Error("Please include  password paramater `-p`.")
		os.Exit(1)
	}

	switch dbType {
	case "mysql":
		runLocalBackupWithMySQL()
	case "postgres", "psql":
		runLocalBackupWithPSQL()
	default:
		utility.Error("UnsuppError: Failed to create database backup: database is closedorted database type: %s. Use 'mysql', or 'postgres'.\n", dbType)
		os.Exit(1)
	}
}

func runLocalBackupWithMySQL() {
	if schedule == "" {
		performLocalBackupWithMySQL(config.Envs.XNAP_DB, schedule)
	} else {
		scheduleLocalBackup(config.Envs.XNAP_DB, schedule)
	}
}

func runLocalBackupWithPSQL() {
	if schedule == "" {
		performLocalBackupWithPSQL(config.Envs.XNAP_DB, schedule)
	} else {
		scheduleLocalBackup(config.Envs.XNAP_DB, schedule)
	}
}

func performLocalBackupWithMySQL(databaseName, _ string) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, databaseName)
	filename = filepath.Base(sourcePath)
	filename = common.GenerateVersionedFilename(versionNum, filename)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		duration = time.Since(start).Seconds()
		database.SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "backup", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err)
		}
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)
	utility.Info("Starting backup for %s file ...", utility.Yellow(filename))

	err = performLocalBackup(sourcePath, backupDirPath)
	if err != nil {
		duration = time.Since(start).Seconds()
		database.SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "backup", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		utility.Error("%v", err.Error())
		os.Exit(1)
	}

	backupRecord := Backup{
		FileName:   filename,
		SourcePath: sourcePath,
		BackupPath: backupDirPath,
	}

	err = db.Create(&backupRecord).Error
	if err != nil {
		duration = time.Since(start).Seconds()
		database.SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "backup", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		utility.Error("Failed to insert backup record: %v", err)
		os.Exit(1)
	}

	duration = time.Since(start).Seconds()
	err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "backup", command, status, errorMessage, dbUser, duration)
	if err != nil {
		utility.Error("Error logging backup command: %v", err)
	}

	utility.Success("Backup metadata inserted successfully into the %s database.", utility.Yellow(databaseName))
}
func performLocalBackupWithPSQL(databaseName, _ string) {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", POSTGRES_DB_HOST, POSTGRES_DB_PORT, dbUser, dbPassword, databaseName)
	filename = filepath.Base(sourcePath)
	filename = common.GenerateVersionedFilename(versionNum, filename)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		duration = time.Since(start).Seconds()
		database.SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "backup", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err)
		}
		utility.Error("Failed to connect to MySQL: %s", err)
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)
	utility.Info("Starting backup for %s file ...", utility.Yellow(filename))

	err = performLocalBackup(sourcePath, backupDirPath)
	if err != nil {
		duration = time.Since(start).Seconds()
		database.SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "backup", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		utility.Error("%v", err.Error())
		os.Exit(1)
	}

	backupRecord := Backup{
		FileName:   filename,
		SourcePath: sourcePath,
		BackupPath: backupDirPath,
	}

	err = db.Create(&backupRecord).Error
	if err != nil {
		duration = time.Since(start).Seconds()
		database.SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "backup", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		utility.Error("Failed to insert backup record: %v", err)
		os.Exit(1)
	}

	duration = time.Since(start).Seconds()
	err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "backup", command, status, errorMessage, dbUser, duration)
	if err != nil {
		utility.Error("Error logging backup command: %v", err)
	}

	utility.Success("Backup metadata inserted successfully into the %s database.", utility.Yellow(databaseName))
}

func scheduleLocalBackup(databaseName, schedule string) {
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

	utility.Info("Backup scheduled for %s\n", utility.Yellow(scheduleTime.String()))

	WG.Add(1)
	// Context to handle cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Wait until the scheduled time
	waitDuration := time.Until(scheduleTime)
	time.AfterFunc(waitDuration, func() {
		defer WG.Done() // Mark the backup task as done
		utility.Info("Starting scheduled backup process at %s\n", time.Now().Format(time.RFC1123))
		if dbType == "mysql" {
			performLocalBackupWithMySQL(databaseName, schedule)
		} else if dbType == "psql" || dbType == "postgres" {
			performLocalBackupWithPSQL(databaseName, schedule)
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
		utility.Info("Backup completed successfully.")
	case <-ctx.Done():
		utility.Info("Backup process was canceled.")
	}
}

func performLocalBackup(sourcePath, backupFolderPath string) error {
	backupFullPath := filepath.Join(backupFolderPath, filename)

	// Ensure the backup folder exists, create it if not
	err := os.MkdirAll(backupFolderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create backup folder: %v", err)
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	backupFile, err := os.Create(backupFullPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer backupFile.Close()

	_, err = io.Copy(backupFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	utility.Success("%s is backup at location %s", utility.Yellow(filename), utility.Yellow(backupFullPath))

	return nil
}
