package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/utils"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func PrintTests(testsData []*models.Test, courseItem *models.CourseItem, printAnswerSheet bool, separateAnswerPage bool) (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	tmpFolder, err := utils.CreateTmpFolder(workDir)
	if err != nil {
		return "", err
	}

	var (
		joiner []string
		mu     sync.Mutex
		wg     sync.WaitGroup
		sem    = make(chan struct{}, 20) // limit to 10 concurrent generations
	)

	for _, testData := range testsData {
		wg.Add(1)

		// acquire semaphore
		sem <- struct{}{}

		go func(testData *models.Test) {
			testData.CourseItem = courseItem
			defer wg.Done()
			defer func() { <-sem }() // release semaphore

			testOutputDir, err := utils.CreateFolder(
				filepath.Join(tmpFolder, strconv.Itoa(int(testData.ID))),
			)
			if err != nil {
				panic(err)
			}

			assetDir, err := utils.CreateFolder(filepath.Join(tmpFolder, "assets"))
			if err != nil {
				panic(err)
			}

			finalTestPath, finalTestPages, err := GenerateTestContent(testData, workDir, assetDir, testOutputDir)
			if err != nil {
				panic(err)
			}

			var paths []string

			if printAnswerSheet {
				if len(testData.Instances) == 0 {
					answerSheetDir, err := utils.CreateFolder(filepath.Join(testOutputDir, "instances"))
					if err != nil {
						panic(err)
					}
					answerSheetPrinter := AnswerSheetPrinter{
						WorkDir:    workDir,
						OutputDir:  answerSheetDir,
						OutputName: "common" + uuid.NewString() + ".pdf",
					}
					answerSheetPath, answerSheetPages := answerSheetPrinter.GenerateAnswerSheets(testData, nil, separateAnswerPage)
					if ((answerSheetPages + finalTestPages) % 2) == 0 {
						paths = append(paths, answerSheetPath, finalTestPath)
					} else {
						paths = append(paths, answerSheetPath, finalTestPath, workDir+"/assets/blank.pdf")
					}
				} else {
					for _, instance := range testData.Instances {
						answerSheetDir, err := utils.CreateFolder(filepath.Join(testOutputDir, "instances"))
						if err != nil {
							panic(err)
						}
						answerSheetPrinter := AnswerSheetPrinter{
							WorkDir:    workDir,
							OutputDir:  answerSheetDir,
							OutputName: strconv.Itoa(int(instance.ID)) + ".pdf",
						}
						answerSheetPath, answerSheetPages := answerSheetPrinter.GenerateAnswerSheets(testData, &instance, separateAnswerPage)
						if ((answerSheetPages + finalTestPages) % 2) == 0 {
							paths = append(paths, answerSheetPath, finalTestPath)
						} else {
							paths = append(paths, answerSheetPath, finalTestPath, workDir+"/assets/blank.pdf")
						}
					}
				}
			} else {
				paths = append(paths, finalTestPath)
			}

			// append to shared slice safely
			mu.Lock()
			joiner = append(joiner, paths...)
			mu.Unlock()

		}(testData)
	}

	wg.Wait()

	mergeDir, err := utils.CreateFolder(
		filepath.Join(tmpFolder, "merging"),
	)
	if err != nil {
		panic(err)
	}
	return MergeFiles(joiner, mergeDir), nil
}

func MergeFiles(filesToJoin []string, tmpPath string) string {
	batchSize := 50 // merge 50 PDFs at a time

	tempFiles := []string{}
	for i := 0; i < len(filesToJoin); i += batchSize {
		end := i + batchSize
		if end > len(filesToJoin) {
			end = len(filesToJoin)
		}

		batch := filesToJoin[i:end]
		tempFile := filepath.Join(
			tmpPath,
			fmt.Sprintf("temp_%d.pdf", i/batchSize),
		)
		err := api.MergeCreateFile(batch, tempFile, false, nil)
		if err != nil {
			log.Fatalf("Failed to merge batch: %v", err)
		}
		tempFiles = append(tempFiles, tempFile)
	}

	utils.DebugPrintJSON(tempFiles)

	if len(tempFiles) == 1 {
		return tempFiles[0]
	} else {
		// Merge all temporary files into the final output
		finalMergeFile := filepath.Join(
			tmpPath,
			"merged_final.pdf",
		)
		err := api.MergeCreateFile(tempFiles, finalMergeFile, false, nil)
		if err != nil {
			log.Fatalf("Failed to merge final PDF: %v", err)
		}
		return finalMergeFile
	}
}
