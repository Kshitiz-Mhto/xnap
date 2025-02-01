package common

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/savioxavier/termlink"
	"golang.org/x/term"
)

var (
	OutputFields string
	// OutputFormat for custom format output
	OutputFormat string
	// RegionSet picks the region to connect to, if you use this option it will use it over the default region
	RegionSet string = "Localhost/127.0.0.1"

	DefaultYes bool
	// PrettySet : Prints the json output in pretty format
	PrettySet bool
	// VersionCli is set from outside using ldflags
	VersionCli = "1.0.0"
	// DateCli is set from outside using ldflags
	DateCli = time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)
)

// IssueMessage is the message to be displayed when an error is returned
func IssueMessage() {
	gitIssueLink := termlink.ColorLink("GitHub issue", "https://github.com/Kshitiz-Mhto/dsync/issues", "green")
	fmt.Printf("Please check if you are using the latest version of CLI and retry the command \nIf you are still facing issues, please report it on our community slack or open a %s \n", gitIssueLink)
}

func EscapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func PromptForPassword() string {
	fmt.Print("Enter password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Print a newline after input
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}
	return string(bytePassword)
}
