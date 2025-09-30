package crons

import (
	"log"
	"os"
	"strings"
	"time"
)

func ClearTempDir() {
	entries, err := os.ReadDir("./temp/")
	if err != nil {
		log.Fatal(err)
	}

	// Reference point: 7 days ago
	weekAgo := time.Now().AddDate(0, 0, -7)

	for _, v := range entries {
		dateTime := strings.Split(v.Name(), ".")[0]

		if len(dateTime) != 15 {
			continue
		}

		t, err := time.Parse("20060102_150405", dateTime)
		if err != nil {
			continue
		}

		// Check if older than a week
		if t.Before(weekAgo) {
			// Example: delete old file
			err := os.RemoveAll("./temp/" + v.Name())
			if err != nil {
				log.Printf("failed to remove %s: %v", v.Name(), err)
			}
		}
	}
}
