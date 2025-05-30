package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Kshitiz-Mhto/xnap/pkg/config"
	"github.com/spf13/cobra"
)

var (
	directory = config.Envs.KEY_LOCATION
	filename  = config.Envs.KEY_FILE
)

var CRYPTOgen = &cobra.Command{
	Use:   "gen",
	Short: "Generate AES key",
	Long:  "Generate AES key with 256 bit key size",
	Run:   setupAESkeyInSys,
}

func setupAESkeyInSys(cmd *cobra.Command, args []string) {

	key, err := GenerateAESKey()
	if err != nil {
		fmt.Printf("Error generating AES key: %v\n", err)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}

	xnapDir := filepath.Join(homeDir, directory)
	err = os.MkdirAll(xnapDir, 0700)
	if err != nil {
		fmt.Printf("Error creating .xnap directory: %v\n", err)
		return
	}

	keysFile := filepath.Join(xnapDir, filename)

	keyHex := hex.EncodeToString(key)

	err = os.WriteFile(keysFile, []byte(keyHex), 0600)
	if err != nil {
		fmt.Printf("Error writing key to file: %v\n", err)
		return
	}

	fmt.Printf("AES key generated successfully and saved to: %s\n", keysFile)
	fmt.Printf("Key size: %d bytes\n", len(key))
}

// GenerateAESKey generates the AES key with 32 bytes size
func GenerateAESKey() ([]byte, error) {
	size, err := strconv.Atoi(config.Envs.AES_KEY_SIZE)
	if err != nil {
		return nil, err
	}

	key := make([]byte, size)
	_, err = rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// LoadAESKey loads the AES key from the keys file
func LoadAESKey() ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting xnap directory: %v", err)
	}

	keysFile := filepath.Join(homeDir, directory, filename)

	if _, err := os.Stat(keysFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("keys file does not exist: %s", keysFile)
	}

	keyHex, err := os.ReadFile(keysFile)
	if err != nil {
		return nil, fmt.Errorf("error reading keys file: %v", err)
	}

	key, err := hex.DecodeString(string(keyHex))
	if err != nil {
		return nil, fmt.Errorf("error decoding key: %v", err)
	}

	return key, nil
}
