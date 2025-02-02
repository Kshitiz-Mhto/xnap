package local

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Kshitiz-Mhto/xnap/pkg/config"
	"github.com/Kshitiz-Mhto/xnap/utility"
	"github.com/spf13/cobra"
)

var (
	MySQL_DB_USER     string = config.Envs.MySQL_DB_USER
	MySQL_DB_PASSWORD string = config.Envs.MySQL_DB_PASSWORD
	MySQL_DB_HOST     string = config.Envs.MySQL_DB_HOST
	MySQL_DB_PORT     string = config.Envs.MySQL_DB_PORT

	POSTGRES_DB_USER     string = config.Envs.POSTGRES_DB_USER
	POSTGRES_DB_PASSWORD string = config.Envs.POSTGRES_DB_PASSWORD
	POSTGRES_DB_HOST     string = config.Envs.POSTGRES_DB_HOST
	POSTGRES_DB_PORT     string = config.Envs.POSTGRES_DB_PORT
)

var (
	dbType        string
	dbUser        string
	dbPassword    string
	backupDirPath string
	sourcePath    string
	restorePath   string
	filename      string
	versionNum    string
	promptPass    bool
	schedule      string
	status        string
	errorMessage  string
	command       string
	start         time.Time
	duration      float64
	WG            sync.WaitGroup
)

var LocalCMD = &cobra.Command{
	Use:   "local",
	Short: "Backup and Restore local files in database",
	Long:  "Backup and Restore local files in database where it copies the local file to backup location and sotre only the reference in the backup table along wih metadata. And Restores files by using the backup file path and copy it to its orginal or custom location",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {

	LocalCMD.AddCommand(LocalBackupCmd)
	LocalCMD.AddCommand(LocalRestoreCmd)

	LocalBackupCmd.Flags().StringVarP(&backupDirPath, "path", "P", ".", "Location of the backup storage")
	LocalBackupCmd.Flags().StringVarP(&sourcePath, "source", "S", "", "Complete file path")
	LocalBackupCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "Schedule backup of database")
	LocalBackupCmd.Flags().StringVarP(&dbUser, "user", "u", "", "Database username")
	LocalBackupCmd.Flags().BoolVarP(&promptPass, "password", "p", false, "Prompt for password (no inline input)")
	LocalBackupCmd.Flags().StringVarP(&dbType, "type", "t", "", "Specify the type of database( Required)")
	LocalBackupCmd.Flags().StringVarP(&versionNum, "version", "v", "", "Specify the backup version of your file")

	LocalBackupCmd.MarkFlagsRequiredTogether("type", "user", "source", "path", "version")
}

func PerformLocalBackup(sourcePath, backupFolderPath string) error {
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
