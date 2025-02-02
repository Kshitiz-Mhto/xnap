package database

import (
	"fmt"

	"github.com/Kshitiz-Mhto/xnap/cli/alert"
	"github.com/Kshitiz-Mhto/xnap/cli/logs"
	"github.com/Kshitiz-Mhto/xnap/pkg/config"
	"github.com/Kshitiz-Mhto/xnap/utility"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LogCommand(dbType, dbUser, dbPassword, host, port, action, command, status, errorMessage string, userName string, duration float64) error {
	var err error

	logEntry := &logs.Log{
		Action:            action,
		Command:           command,
		Status:            status,
		ErrorMessage:      errorMessage,
		UserName:          userName,
		ExecutionDuration: duration,
	}

	switch dbType {
	case "mysql":
		err = AddLogToMysql(dbUser, dbPassword, host, port, logEntry)
	case "psql", "postgres":
		err = AddLogtoPSQL(dbUser, dbPassword, host, port, logEntry)
	}

	return err
}

func AddLogToMysql(dbUser, dbPassword, host, port string, logEntry *logs.Log) error {
	var lastLogEntry logs.Log

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, host, port, config.Envs.XNAP_DB)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to MySQL: %s", err)
		return err
	}

	defer utility.CloseDBConnection(db)

	if err := db.Create(logEntry).Error; err != nil {
		return err
	}

	if err := db.Order("created_at DESC").First(&lastLogEntry).Error; err != nil {
		fmt.Println("Error fetching the last row:", err)
	}

	if logEntry.Status == config.Envs.BACKUP_OR_RESTORE_STATUS {
		vars := map[string]interface{}{
			"ID":                lastLogEntry.ID,
			"Action":            lastLogEntry.Action,
			"Command":           lastLogEntry.Command,
			"Status":            lastLogEntry.Status,
			"ErrorMessage":      lastLogEntry.ErrorMessage,
			"UserName":          lastLogEntry.UserName,
			"ExecutionDuration": lastLogEntry.ExecutionDuration,
			"CreatedAt":         lastLogEntry.CreatedAt,
			"UpdatedAt":         lastLogEntry.UpdatedAt,
			"dbType":            dbType,
		}
		alert.HTMLTemplateEmailHandler(config.Envs.OWNER_EMAIL, vars)
	}
	utility.Success("Log is enteried successfully!!")

	return nil
}

func AddLogtoPSQL(dbUser, dbPassword, host, port string, logEntry *logs.Log) error {
	var lastLogEntry logs.Log

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, dbUser, dbPassword, config.Envs.XNAP_DB)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		utility.Error("Failed to connect to PSQL: %s", err)
		return err
	}

	defer utility.CloseDBConnection(db)

	if err := db.Create(logEntry).Error; err != nil {
		return err
	}

	if err := db.Order("created_at DESC").First(&lastLogEntry).Error; err != nil {
		return err
	}

	if logEntry.Status == config.Envs.BACKUP_OR_RESTORE_STATUS {
		vars := map[string]interface{}{
			"ID":                lastLogEntry.ID,
			"Action":            lastLogEntry.Action,
			"Command":           lastLogEntry.Command,
			"Status":            lastLogEntry.Status,
			"ErrorMessage":      lastLogEntry.ErrorMessage,
			"UserName":          lastLogEntry.UserName,
			"ExecutionDuration": lastLogEntry.ExecutionDuration,
			"CreatedAt":         lastLogEntry.CreatedAt,
			"UpdatedAt":         lastLogEntry.UpdatedAt,
			"dbType":            dbType,
		}

		alert.HTMLTemplateEmailHandler(config.Envs.OWNER_EMAIL, vars)
	}

	utility.Success("Log is enteried successfully!!")

	return nil
}

func SetFailureStatus(msg string) {
	status = "failure"
	errorMessage = msg
}
