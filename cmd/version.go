package cmd

import (
	"fmt"
	"runtime"

	"github.com/Kshitiz-Mhto/dsync/utility/common"
	"github.com/spf13/cobra"
)

const logo = `
   ______   _______           _        _______ 
  (  __  \ (  ____ \|\     /|( (    /|(  ____ \
  | (  \  )| (    \/( \   / )|  \  ( || (    \/
  | |   ) || (_____  \ (_) / |   \ | || |      
  | |   | |(_____  )  \   /  | (\ \) || |      
  | |   ) |      ) |   ) (   | | \   || |      
  | (__/  )/\____) |   | |   | )  \  || (____/\
  (______/ \_______)   \_/   |/    )_)(_______/                                               

`

var (
	quiet      bool
	verbose    bool
	versionCMD = &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Run: func(cmd *cobra.Command, args []string) {
			switch {
			case verbose:
				fmt.Print(logo)
				fmt.Printf("Client version: v%s\n", common.VersionCli)
				fmt.Printf("Go version (client): %s\n", runtime.Version())
				fmt.Printf("Build date (client): %s\n", common.DateCli)
				fmt.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)
				// common.CheckVersionUpdate()
			case quiet:
				fmt.Printf("v%s\n", common.VersionCli)
			default:
				// common.CheckVersionUpdate()
				fmt.Printf("dSync CLI v%s\n", common.VersionCli)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCMD)
	versionCMD.Flags().BoolVarP(&quiet, "quiet", "q", false, "Use quiet output for simple output")
	versionCMD.Flags().BoolVarP(&verbose, "verbose", "v", false, "Use verbose output to see full information")
}
