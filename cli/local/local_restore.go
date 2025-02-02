package local

import (
	"os"
	"strings"
	"time"

	"github.com/Kshitiz-Mhto/xnap/utility"
	"github.com/Kshitiz-Mhto/xnap/utility/common"
	"github.com/spf13/cobra"
)

var LocalRestoreCmd = &cobra.Command{
	Use:     "restore",
	Aliases: []string{"reset", "restores"},
	Short:   "Restores files to its original path",
	Long:    "Restores files by using the backup file path and copy it to its orginal or custom location",
	Example: "xnap local restore --type <database_type> -u <username> -p --path <to/path>  --schedule <schedule_HH:MM>",
	Run:     runLocalRestoreCommand,
}

func runLocalRestoreCommand(cmd *cobra.Command, args []string) {
	restorePath, _ = cmd.Flags().GetString("path")
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

}
