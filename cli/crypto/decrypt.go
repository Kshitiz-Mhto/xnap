package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kshitiz-Mhto/xnap/utility"
	"github.com/spf13/cobra"
)

var (
	backupFileNamePath string
	decryptedFilePath  string
)

var CRYPTOdecrypt = &cobra.Command{
	Use:     "decrypt",
	Aliases: []string{"de"},
	Short:   "Decrypt the dump file",
	Example: "xnap crypto decrypt --path <filepath> --file <encrypted-file>",
	Run:     runDecryptionOfDumpFile,
}

func runDecryptionOfDumpFile(cmd *cobra.Command, args []string) {

	backupFileNamePath, _ = cmd.Flags().GetString("file")
	decryptedFilePath, _ = cmd.Flags().GetString("path")

	// Validate that encrypted file path is provided
	if backupFileNamePath == "" {
		utility.Error("Error: Please provide the encrypted file path using --file flag")
		return
	}

	// Resolve the full path to the encrypted file
	var fullEncryptedPath string
	var err error

	// The --file flag contains the full path to encrypted file
	fullEncryptedPath, err = filepath.Abs(backupFileNamePath)
	if err != nil {
		utility.Error("Error resolving encrypted file path: %v\n", err)
		return
	}

	// Check if encrypted file exists
	if _, err := os.Stat(fullEncryptedPath); os.IsNotExist(err) {
		utility.Error("Error: Encrypted file does not exist: %s\n", fullEncryptedPath)
		return
	}

	// Decrypt the file
	utility.Error("Decrypting file: %s\n", fullEncryptedPath)

	decryptedData, err := decryptFile(fullEncryptedPath)
	if err != nil {
		utility.Error("Error decrypting file: %v\n", err)
		return
	}

	// Generate output file name (remove .enc extension if present)
	baseFileName := filepath.Base(fullEncryptedPath)
	outputFileName := strings.TrimSuffix(baseFileName, ".xnap")

	// Create full output path using the decryptedFilePath (--path flag)
	outputFilePath := filepath.Join(decryptedFilePath, outputFileName)
	outputFilePath, err = filepath.Abs(outputFilePath)
	if err != nil {
		utility.Error("Error resolving output file path: %v\n", err)
		return
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputFilePath)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		utility.Error("Error creating output directory: %v\n", err)
		return
	}

	// Write decrypted data to output file
	err = os.WriteFile(outputFilePath, decryptedData, 0600)
	if err != nil {
		utility.Error("Error writing decrypted file: %v\n", err)
		return
	}

	utility.Success("Decryption completed successfully!\n")
	utility.Success("Decrypted file saved at: %s\n", utility.Green(outputFilePath))
	utility.Success("Original encrypted file: %s\n", utility.Yellow(fullEncryptedPath))
	utility.Success("File size: %d bytes\n", len(decryptedData))
}

// decryptFile decrypts an encrypted backup file
func decryptFile(encryptedFilePath string) ([]byte, error) {
	// Load the AES key from the keys file
	key, err := LoadAESKey()
	if err != nil {
		return nil, fmt.Errorf("failed to load AES key: %v", err)
	}

	// Read the encrypted file
	ciphertext, err := os.ReadFile(encryptedFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted file: %v", err)
	}

	// Validate file size
	if len(ciphertext) == 0 {
		return nil, fmt.Errorf("encrypted file is empty")
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	// Extract nonce and encrypted data
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short: expected at least %d bytes, got %d", nonceSize, len(ciphertext))
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data (wrong key or corrupted file): %v", err)
	}

	return plaintext, nil
}

func init() {
	CRYPTOdecrypt.Flags().StringVarP(&backupFileNamePath, "file", "f", "", "File path for encrypted backup file (Required)")
	CRYPTOdecrypt.Flags().StringVarP(&decryptedFilePath, "path", "P", ".", "File path to save the decrypted backup file (default: current directory)")

	CRYPTOdecrypt.MarkFlagRequired("file")
}
