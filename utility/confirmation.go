package utility

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"gorm.io/gorm"
)

// retrieveUserInput is a function that can retrieve user input in form of string. By default,
// it will prompt the user. In test, you can replace this with code that returns the appropriate response.
var retrieveUserInput = func(message string) (string, error) {
	return readUserInput(os.Stdin, message)
}

// readUserInput is a io.Reader to read user input from. It is meant to allow simplified testing
// as to-be-read inputs can be injected conveniently.
func readUserInput(in io.Reader, message string) (string, error) {
	reader := bufio.NewReader(in)
	YellowConfirm("Are you sure you want to %s (y/N) ? ", message)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	answer = strings.TrimRight(answer, "\r\n")

	return strings.ToLower(answer), nil
}

// AskForConfirm parses and verifies user input for confirmation.
func AskForConfirm(message string) error {
	answer, err := retrieveUserInput(message)
	if err != nil {
		Error("Unable to parse users input: %s", err)
	}

	if answer != "y" && answer != "ye" && answer != "yes" {
		return fmt.Errorf("invalid user input")
	}

	return nil
}

// UserConfirmedDeletion builds a message to ask the user to confirm delete
// a resource and then sends it through to AskForConfirm to
// parses and verifies user input.
func UserConfirmedDeletion(resourceType string, ignoringConfirmed bool, objectToDelete string) bool {
	if !ignoringConfirmed {
		message := fmt.Sprintf("delete the %s %s", Green(objectToDelete), resourceType)
		err := AskForConfirm(message)
		if err != nil {
			return false
		}
	}

	return true
}

// UserAccepts is a function that can retrieve user input in form of string and checks if it is a yes.
func UserAccepts(in io.Reader) (bool, error) {
	reader := bufio.NewReader(in)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	answer = strings.TrimRight(answer, "\r\n")
	strings.ToLower(answer)
	if answer != "y" && answer != "ye" && answer != "yes" {
		return false, fmt.Errorf("invalid user input")
	}

	return true, nil
}

func CloseDBConnection(db *gorm.DB) {
	// Ensure the connection is closed when the function exits
	DB, err := db.DB() // Get the underlying *sql.DB instance
	if err != nil {
		Error("Failed to get SQL DB instance: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := DB.Close(); err != nil {
			Error("Failed to close database connection: %v", err)
		}
	}()
}
