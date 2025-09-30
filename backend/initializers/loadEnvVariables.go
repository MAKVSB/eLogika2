package initializers

import (
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var GlobalAppConfig *AppConfig

type AppConfig struct {
	PORT                     int64
	MODE                     string
	PROTOCOL                 string
	CERTPATH                 *string
	CERTNAME                 *string
	PLATFORM                 string
	GIN_RELEASE_MODE         bool
	DB_URL                   string
	ACCESS_SECRET            []byte
	ACCESS_LENGTH            time.Duration
	REFRESH_SECRET           []byte
	REFRESH_LENGTH           time.Duration
	API_SECRET               []byte
	API_LENGTH               time.Duration
	ACCESS_TOKEN_REVOKE_SYNC bool
	UPLOADS_DESTINATION      string
	INBUS_BASE_URL           string
	INBUS_CLIENT_ID          string
	INBUS_CLIENT_SECRET      string
}

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, continuing with system env vars")
	}

	accessLen, err := parseAdvancedDuration(getEnv("ACCESS_LENGTH", "15m"))
	if err != nil {
		log.Fatalf("Invalid ACCESS_LENGTH format: %v", err)
	}
	refreshLen, err := parseAdvancedDuration(getEnv("REFRESH_LENGTH", "7d"))
	if err != nil {
		log.Fatalf("Invalid REFRESH_LENGTH format: %v", err)
	}
	apiLen, err := parseAdvancedDuration(getEnv("API_LENGTH", "30d"))
	if err != nil {
		log.Fatalf("Invalid REFRESH_LENGTH format: %v", err)
	}

	GlobalAppConfig = &AppConfig{
		PORT:                     getEnvInt("PORT", 8080),
		MODE:                     getEnv("MODE", "prod"),
		PROTOCOL:                 getEnv("PROTOCOL", "https"),
		CERTPATH:                 getEnvNilable("CERTPATH"),
		CERTNAME:                 getEnvNilable("CERTNAME"),
		PLATFORM:                 getEnvPlatform("PLATFORM"),
		GIN_RELEASE_MODE:         getEnvBool("GIN_RELEASE_MODE", false),
		DB_URL:                   getEnv("DB_URL", ""),
		ACCESS_SECRET:            []byte(getEnv("ACCESS_SECRET", "")),
		ACCESS_LENGTH:            accessLen,
		REFRESH_SECRET:           []byte(getEnv("REFRESH_SECRET", "")),
		REFRESH_LENGTH:           refreshLen,
		API_SECRET:               []byte(getEnv("REFRESH_SECRET", "")),
		API_LENGTH:               apiLen,
		ACCESS_TOKEN_REVOKE_SYNC: getEnvBool("ACCESS_TOKEN_REVOKE_SYNC", false),
		UPLOADS_DESTINATION:      getEnv("UPLOADS_DESTINATION", "./uploads"),
		INBUS_BASE_URL:           getEnv("INBUS_BASE_URL", "https://inbus.vsb.cz/"),
		INBUS_CLIENT_ID:          getEnv("INBUS_CLIENT_ID", ""),
		INBUS_CLIENT_SECRET:      getEnv("INBUS_CLIENT_SECRET", ""),
	}
}

func getEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvNilable(key string) *string {
	val := os.Getenv(key)
	if val == "" {
		return nil
	}
	return &val
}

func getEnvInt(key string, fallback int64) int64 {
	val := os.Getenv(key)
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return i
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	val := strings.ToLower(os.Getenv(key))
	if val == "true" || val == "1" {
		return true
	}
	if val == "false" || val == "0" {
		return false
	}
	return fallback
}

func getEnvPlatform(key string) string {
	val := os.Getenv(key)
	if val == "" || val == "AUTO" {
		return runtime.GOOS
	}
	return val
}

func parseAdvancedDuration(input string) (time.Duration, error) {
	re := regexp.MustCompile(`(?i)(\d+)([smhdw])`)
	matches := re.FindAllStringSubmatch(input, -1)

	var total time.Duration
	unitMap := map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": time.Hour * 24,
		"w": time.Hour * 24 * 7,
	}

	for _, match := range matches {
		val, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}
		unit := strings.ToLower(match[2])
		multiplier, ok := unitMap[unit]
		if !ok {
			return 0, err
		}
		total += time.Duration(val) * multiplier
	}

	return total, nil
}
