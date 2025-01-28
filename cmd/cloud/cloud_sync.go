package cloud

import (
	"github.com/spf13/cobra"
)

var cloudSyncCmd = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"up"},
	Short:   "File-to-Cloud Synchronization",
	Long:    "uploading files from a local system to a cloud storage provider",
	Example: "dsync cloud sync --type <cloud_storage_type",
	Run:     operationUploadFileToCloud,
}

func operationUploadFileToCloud(cmd *cobra.Command, args []string) {
	panic("unimplemented")
}
