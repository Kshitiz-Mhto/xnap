/*
Copyright © 2025 Kshitiz Mhto <kshitizmhto101@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cli

import (
	"fmt"
	"os"

	"github.com/Kshitiz-Mhto/xnap/cli/cloud"
	"github.com/Kshitiz-Mhto/xnap/cli/crypto"
	"github.com/Kshitiz-Mhto/xnap/cli/database"
	"github.com/Kshitiz-Mhto/xnap/cli/local"
	"github.com/spf13/cobra"
)

var version bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xnap",
	Short: "CLI for database and storage files backup and restore utility",
	Long:  "xnap is a CLI tool designed to simplify the backup and restoration including databases, and local system files.",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			versionCMD.Run(cmd, args)
		} else {
			fmt.Print(logo)
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.

	rootCmd.AddCommand(database.DBCmd)
	rootCmd.AddCommand(cloud.CloudSyncCmd)
	rootCmd.AddCommand(local.LocalCMD)
	rootCmd.AddCommand(crypto.CRYPTOgen)
	rootCmd.AddCommand(crypto.CRYPTOdecrypt)

	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Print the version of the CLI")
}
