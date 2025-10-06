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

func PrintTests(testsData []*models.Test, courseItem *models.CourseItem, printAnswerSheet bool, separateAnswerPage bool) string {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	tmpFolder := utils.CreateTmpFolder(workDir)

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

			testOutputDir := utils.CreateFolder(
				filepath.Join(tmpFolder, strconv.Itoa(int(testData.ID))),
			)

			testPrinter := TestPrinter{
				WorkDir:   workDir,
				AssetDir:  utils.CreateFolder(filepath.Join(tmpFolder, "assets")),
				OutputDir: testOutputDir,
			}

			finalTestPath, finalTestPages := testPrinter.GenerateTestContent(testData)

			var paths []string

			if printAnswerSheet {
				if len(testData.Instances) == 0 {
					answerSheetDir := utils.CreateFolder(filepath.Join(testOutputDir, "instances"))
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
						answerSheetDir := utils.CreateFolder(filepath.Join(testOutputDir, "instances"))
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

	mergeDir := utils.CreateFolder(
		filepath.Join(tmpFolder, "merging"),
	)
	return MergeFiles(joiner, mergeDir)
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
