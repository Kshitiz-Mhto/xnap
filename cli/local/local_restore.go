package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

var LocalRestoreCmd = &cobra.Command{
	Use:     "restore",
	Aliases: []string{"reset", "restores"},
	Short:   "Restores files to its original path",
	Long:    "Restores files by using the backup file path and copy it to its orginal or custom location",
	Example: "xnap local restore --type <database_type> -u <username> -p --file <filename> --path <to/path/restore_location> --version <version_number> --schedule <schedule_HH:MM>",
	Run:     runLocalRestoreCommand,
}

func runLocalRestoreCommand(cmd *cobra.Command, args []string) {
	restorePath, _ = cmd.Flags().GetString("path")
	filename, _ = cmd.Flags().GetString("file")
	versionNum, _ = cmd.Flags().GetString("version")
	restorePath, _ = filepath.Abs(restorePath)
	filenamePattern1 = regexp.MustCompile(`^\.(.*)_v(\d+)$`)         // .filename_v{number}
	filenamePattern2 = regexp.MustCompile(`^(.*)_v(\d+)(\.[^.]*)?$`) // filename_v{number}.ext
	command = strings.Join(os.Args, " ")
	start = time.Now()
	status = "success"
	errorMessage = ""

	err := os.MkdirAll(restorePath, os.ModePerm)
	if err != nil {
		utility.Error("Error resolving backup file path: %v", err)
		os.Exit(1)
	}
	if err != nil {
		utility.Error("Failed to create restore path for file.")
		os.Exit(1)
	}

	if promptPass {
		dbPassword = common.PromptForPassword()
	} else {
		utility.Error("Please include  password paramater `-p`.")
		os.Exit(1)
	}

	switch dbType {
	case "mysql":
		runLocalRestoreWithMySQL()
	case "postgres", "psql":
		runLocalRestoreWithPSQL()
	default:
		utility.Error("UnsuppError: Failed to create database backup: database is closedorted database type: %s. Use 'mysql', or 'postgres'.\n", dbType)
		os.Exit(1)
	}

}

func runLocalRestoreWithMySQL() {
	if schedule == "" {
		performLocalRestoreWithMySQL(config.Envs.XNAP_DB, schedule)
	} else {
		scheduleLocalRestore(config.Envs.XNAP_DB, schedule)
	}
}

func runLocalRestoreWithPSQL() {
	if schedule == "" {
		performLocalRestoreWithPSQL(config.Envs.XNAP_DB, schedule)
	} else {
		scheduleLocalRestore(config.Envs.XNAP_DB, schedule)
	}
}

func performLocalRestoreWithMySQL(databaseName, _ string) {
	var matchedBackup *Backup
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, databaseName)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		duration = time.Since(start).Seconds()
		SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err)
		}
		os.Exit(1)
	}
	defer utility.CloseDBConnection(db)
	utility.Info("Starting Restoration for %s file from %s database ...", utility.Yellow(filename), utility.Yellow(dbType))

	var listOfBackups []Backup

	err = db.Raw("SELECT * FROM backups WHERE og_file_name = ?", filename).Scan(&listOfBackups).Error

	if err != nil {
		utility.Error("Failed to fetch databases: %v", err)
		duration = time.Since(start).Seconds()
		SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		os.Exit(1)
	}

	if listOfBackups == nil {
		utility.Warning("No backups found for the given source path.")
		duration = time.Since(start).Seconds()
		SetFailureStatus("No backups found for the given source path.")
		err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		os.Exit(1)
	}

	for _, backup := range listOfBackups {
		fileName := backup.FileName

		if matches := filenamePattern1.FindStringSubmatch(fileName); matches != nil {
			if CheckVersionMatch(matches[2], versionNum) {
				matchedBackup = &backup
				break
			}
		} else if matches := filenamePattern2.FindStringSubmatch(fileName); matches != nil {
			if CheckVersionMatch(matches[2], versionNum) {
				matchedBackup = &backup
				break
			}
		}
	}

	if matchedBackup == nil {
		duration = time.Since(start).Seconds()
		SetFailureStatus("Backup object is nil")
		err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		os.Exit(1)
	}

	fullBackupFilePath, err := filepath.Abs(filepath.Join(matchedBackup.BackupPath, matchedBackup.FileName))
	if err != nil {
		utility.Error("Error resolving backup file path: %v", err)
		duration = time.Since(start).Seconds()
		SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		os.Exit(1)
	}
	err = PerformLocalBackup(fullBackupFilePath, restorePath)
	if err != nil {
		utility.Error("%v", err.Error())
		duration = time.Since(start).Seconds()
		SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		os.Exit(1)
	}

	duration = time.Since(start).Seconds()
	err = database.LogCommand(dbType, dbUser, dbPassword, MySQL_DB_HOST, MySQL_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
	if err != nil {
		utility.Error("Error logging backup command: %v", err)
	}

	utility.Success("Restoration successfully into from %s database.", utility.Yellow(databaseName))
}

func performLocalRestoreWithPSQL(databaseName, _ string) {
	var matchedBackup *Backup
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", POSTGRES_DB_HOST, POSTGRES_DB_PORT, dbUser, dbPassword, databaseName)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		duration = time.Since(start).Seconds()
		SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err)
		}
		os.Exit(1)
	}

	defer utility.CloseDBConnection(db)
	utility.Info("Starting Restoration process for %s file from %s database ...", utility.Yellow(filename), utility.Yellow(dbType))

	var listOfBackups []Backup

	err = db.Raw("SELECT * FROM backups WHERE og_file_name = ?", filename).Scan(&listOfBackups).Error

	if err != nil {
		utility.Error("Failed to fetch databases: %v", err)
		duration = time.Since(start).Seconds()
		SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		os.Exit(1)
	}

	if listOfBackups == nil {
		utility.Warning("No backups found for the given source path.")
		duration = time.Since(start).Seconds()
		SetFailureStatus("No backups found for the given source path.")
		err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err)
		}
		os.Exit(1)
	}

	for _, backup := range listOfBackups {
		fileName := backup.FileName

		if matches := filenamePattern1.FindStringSubmatch(fileName); matches != nil {
			if CheckVersionMatch(matches[2], versionNum) {
				matchedBackup = &backup
				break
			}
		} else if matches := filenamePattern2.FindStringSubmatch(fileName); matches != nil {
			if CheckVersionMatch(matches[2], versionNum) {
				matchedBackup = &backup
				break
			}
		}
	}

	if matchedBackup == nil {
		duration = time.Since(start).Seconds()
		SetFailureStatus("Backup object is nil")
		err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		os.Exit(1)
	}

	fullBackupFilePath, err := filepath.Abs(filepath.Join(matchedBackup.BackupPath, matchedBackup.FileName))
	if err != nil {
		utility.Error("Error resolving backup file path: %v", err)
		duration = time.Since(start).Seconds()
		SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		os.Exit(1)
	}

	err = PerformLocalBackup(fullBackupFilePath, restorePath)
	if err != nil {
		utility.Error("%v", err.Error())
		duration = time.Since(start).Seconds()
		SetFailureStatus(err.Error())
		err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
		if err != nil {
			utility.Error("Error logging backup command: %v", err.Error())
		}
		os.Exit(1)
	}

	duration = time.Since(start).Seconds()
	err = database.LogCommand(dbType, dbUser, dbPassword, POSTGRES_DB_HOST, POSTGRES_DB_PORT, "restore", command, status, errorMessage, dbUser, duration)
	if err != nil {
		utility.Error("Error logging backup command: %v", err)
	}

	utility.Success("Restoration successfully into from %s database.", utility.Yellow(databaseName))
}

func scheduleLocalRestore(databaseName, schedule string) {
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
			performLocalRestoreWithMySQL(databaseName, schedule)
		} else if dbType == "psql" || dbType == "postgres" {
			performLocalRestoreWithPSQL(databaseName, schedule)
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
