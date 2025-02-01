package cloud

import "github.com/spf13/cobra"

var cloudRestoreCmd = &cobra.Command{
	Use:     "restore",
	Aliases: []string{"backup"},
	Short:   "Restore your local file from cloud storage",
	Example: "dsyc cloud restore --type <cloud_storage_type>",
	Run:     runCloudRestoration,
}

func runCloudRestoration(cmd *cobra.Command, args []string) {
	panic("implemented")
}
