package cloud

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	cloudType string
	fileName  string
)

var CloudSyncCmd = &cobra.Command{
	Use:   "cloud-sync",
	Short: "File-to-Cloud Synchronization",
	Long:  "Manages operations related to cloud file Synchronization",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	CloudSyncCmd.AddCommand(cloudSyncCmd)
	cloudSyncCmd.AddCommand(cloudRestoreCmd)

	cloudSyncCmd.Flags().StringVarP(&cloudType, "type", "t", "", "Specify the type of cloud service provider. Use `gcp` or `aws` or `civo`")
	cloudSyncCmd.Flags().StringVarP(&fileName, "file", "f", "", "Specify the file path")

	cloudRestoreCmd.Flags().StringVarP(&cloudType, "type", "t", "", "Specify the type of cloud service provider. Use `gcp` or `aws` or `civo`")

	cloudSyncCmd.MarkFlagRequired("type")

}
