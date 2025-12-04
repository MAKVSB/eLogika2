package utils

import (
	"archive/zip"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func generateRandomHex(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CreateTmpFolder(root string) (string, error) {
	// current timestamp with millisecond precision
	timestamp := time.Now().Format("20060102_150405.000") // YYYYMMDD_HHMMSS.mmm

	// random value to avoid collisions
	rnd, err := generateRandomHex(4) // 8 hex chars
	if err != nil {
		return "", fmt.Errorf("failed to generate random hex: %w", err)
	}

	// folder name
	folderName := fmt.Sprintf(root+"/temp/%s_%s", timestamp, rnd)
	return CreateFolder(folderName)
}

func CreateFolder(path string) (string, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create folder %s: %w", path, err)
	}
	return path, nil
}

func CopyFile(source string, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", source, err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destination, err)
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy content from %s to %s: %w", source, destination, err)
	}

	if err := destinationFile.Close(); err != nil {
		return fmt.Errorf("failed to close destination file %s after copying: %w", destination, err)
	}

	return nil
}

func LinkFile(source string, destination string) error {
	if _, err := os.Stat(destination); err == nil {
		return nil
	}

	if err := os.Link(source, destination); err != nil {
		return fmt.Errorf("failed to create hard link from %s to %s: %w", source, destination, err)
	}

	return nil
}

func DownloadFile(url, destination string) (string, error) {
	// Compute MD5 hash of the URL to use as filename
	hash := md5.Sum([]byte(url))
	filename := hex.EncodeToString(hash[:]) + filepath.Ext(url) // keep extension if any
	filepath := filepath.Join(destination, filename)

	// Check if file already exists
	if _, err := os.Stat(filepath); err == nil {
		// File already downloaded
		return filepath, nil
	}

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download: %s", resp.Status)
	}

	// Create local file
	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write data to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func GenerateFileName(dbRef *gorm.DB, ext string) (string, error) {
	var newFileName string
	for i := 1; i > 0; i++ {
		candidate := uuid.New().String() + ext

		var file models.File
		if err := dbRef.Model(&models.File{}).
			Where("stored_name = ?", candidate).
			Limit(1).
			Find(&file).Error; err != nil {
			return "", errors.New("database error during UUID check")
		}

		if file.ID == 0 {
			newFileName = candidate
			break
		}

	}

	return newFileName, nil
}

func ZipFolder(folderPath string) (*os.File, error) {
	_, err := CreateFolder(folderPath + "-bac")
	if err != nil {
		return nil, err
	}

	tmpZip, err := os.CreateTemp(folderPath+"-bac", "bac.zip")
	if err != nil {
		return nil, err
	}

	zipWriter := zip.NewWriter(tmpZip)

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // skip directories
		}

		relPath, err := filepath.Rel(folderPath, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		zipFileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFileWriter, file)
		return err
	})
	if err != nil {
		zipWriter.Close()
		tmpZip.Close()
		return nil, err
	}

	zipWriter.Close()
	tmpZip.Seek(0, 0)
	return tmpZip, nil
}

func ZipFolderError(folderPath string) (*common.ErrorFile, *common.ErrorResponse) {
	zipFile, err := ZipFolder(folderPath)
	if err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to print questions and create error report",
			Details: err.Error(),
		}
	}
	defer zipFile.Close()

	zipBytes, err := io.ReadAll(zipFile)
	if err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to load print error logs",
			Details: err.Error(),
		}
	}

	return &common.ErrorFile{
		Content:  base64.StdEncoding.EncodeToString(zipBytes),
		MimeType: "application/zip",
	}, nil
}
