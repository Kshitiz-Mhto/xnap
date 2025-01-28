package database

import (
	"os"

	"github.com/Kshitiz-Mhto/dsync/utility"
	"github.com/spf13/cobra"
)

var dbSyncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"sy"},
	Short:   "Synchronize the database",
	Long:    "Synchronize master and slave homogeneous or hetrogeneous database",
	Example: "dsync db sync --master <master_database_name> --slave <slave_database_name> --type <database_type>",
	Run:     runDatabseSynchronization,
}

func runDatabseSynchronization(cmd *cobra.Command, args []string) {
	switch dbType {
	case "mysql":
		syncMySQLDatabases()
	case "postgres", "psql":
		syncPostgresDatabases()

	default:
		utility.Error("InvalidationError: Invalid database types provided. Use 'mysql', or 'postgres'.\n")
		os.Exit(1)
	}
}

func syncMySQLDatabases() {
	panic("unimplemented")
}

func syncPostgresDatabases() {
	utility.Warning("Feature yet not implemented!!")
	os.Exit(1)
}
