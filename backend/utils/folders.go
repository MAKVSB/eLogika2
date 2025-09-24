package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

func generateRandomHex(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CreateTmpFolder(root string) string {
	// current timestamp with millisecond precision
	timestamp := time.Now().Format("20060102_150405.000") // YYYYMMDD_HHMMSS.mmm

	// random value to avoid collisions
	rnd, err := generateRandomHex(4) // 8 hex chars
	if err != nil {
		panic(err)
	}

	// folder name
	folderName := fmt.Sprintf(root+"/temp/%s_%s", timestamp, rnd)

	// create folder
	err = os.MkdirAll(folderName, 0755)
	if err != nil {
		panic(err)
	}
	return folderName
}

func CreateFolder(path string) string {
	// folder name
	folderName := fmt.Sprintf(path)

	// create folder
	err := os.MkdirAll(folderName, 0755)
	if err != nil {
		panic(err)
	}
	return folderName
}
