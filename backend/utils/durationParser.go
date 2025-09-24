package utils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseDurationString(s string) (time.Duration, error) {
	re := regexp.MustCompile(`(?i)(\d+)([smhdw])`)
	matches := re.FindAllStringSubmatch(s, -1)

	var total time.Duration
	unitMap := map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": time.Hour * 24,
		"w": time.Hour * 24 * 7,
	}

	for _, match := range matches {
		value, _ := strconv.Atoi(match[1])
		unit := strings.ToLower(match[2])
		total += time.Duration(value) * unitMap[unit]
	}

	return total, nil
}
