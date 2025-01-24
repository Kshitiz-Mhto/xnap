package utility

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

// Green is the function to convert str to green in console
func Green(value string) string {
	newColor := color.FgGreen.Render
	return newColor(value)
}

// Yellow is the function to convert str to yellow in console
func Yellow(value string) string {
	newColor := color.New(color.FgYellow).Render
	return newColor(value)
}

// Red is the function to convert str to red in console
func Red(value string) string {
	newColor := color.New(color.FgRed).Render
	return newColor(value)
}

// Error is the function to handler all error in the Cli
func Error(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Red.Sprintf("Error"), fmt.Sprintf(msg, args...))
}

// Info is the function to handler all info messages in the Cli
func Info(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Blue.Sprintf("Info"), fmt.Sprintf(msg, args...))
}

// Warning is the function to handler all warnings in the Cli
func Warning(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Yellow.Sprintf("Warning"), fmt.Sprintf(msg, args...))
}

func Success(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Green.Sprintf("sucess"), fmt.Sprintf(msg, args...))
}

// YellowConfirm is the function to handler all delete confirm
func YellowConfirm(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s", color.Warn.Sprintf("Warning"), fmt.Sprintf(msg, args...))
}
