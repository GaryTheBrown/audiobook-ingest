package config

import (
	"os"
	"strconv"
	"strings"
)

var (
	ImportDirs        []string
	ExportDir         string
	Port              string
	NoDelete          bool
	HideAPIWarning    bool
	DisableAutoIngest bool
	MultiSwitchCount  int

	// Workers Concurrency Caps
	DetectingWorkers  int
	ConversionWorkers int
	VerifyWorkers     int
	ExportWorkers     int

	// APIS
	ChaptarrURL          string
	ChaptarrAPIKey       string
	AudiobookshelfURL    string
	AudiobookshelfAPIKey string
	KomgaURL             string
	KomgaAPIKey          string
	KavitaURL            string
	KavitaAPIKey         string
)

func init() {
	rawImportDirs := getEnv("IMPORT_DIR", "/import")
	var parsedDirs []string

	for _, dir := range strings.Split(rawImportDirs, ",") {
		trimmed := strings.TrimSpace(dir)
		if trimmed != "" {
			parsedDirs = append(parsedDirs, trimmed)
		}
	}

	ImportDirs = parsedDirs
	ExportDir = getEnv("EXPORT_DIR", "/export")
	Port = getEnv("WEB_PORT", "8080")
	_, NoDelete = os.LookupEnv("NO_DELETE")
	_, HideAPIWarning = os.LookupEnv("HIDE_API_WARNING")
	_, DisableAutoIngest = os.LookupEnv("DISABLE_AUTO_INGEST")
	MultiSwitchCount = parseIntEnv("MULTI_SWITCH_COUNT", 5)

	// Workers Concurrency Allocation with Sane Defaults
	DetectingWorkers = parseIntEnv("DETECTING_WORKERS", 2)
	ConversionWorkers = parseIntEnv("CONVERSION_WORKERS", 1)
	VerifyWorkers = parseIntEnv("VERIFY_WORKERS", 2)
	ExportWorkers = parseIntEnv("EXPORT_WORKERS", 2)

	// APIS
	ChaptarrURL = getEnv("CHAPTARR_URL", "")
	ChaptarrAPIKey = getEnv("CHAPTARR_API_KEY", "")
	AudiobookshelfURL = getEnv("AUDIOBOOKSHELF_URL", "")
	AudiobookshelfAPIKey = getEnv("AUDIOBOOKSHELF_API_KEY", "")
	KomgaURL = getEnv("KOMGA_URL", "")
	KomgaAPIKey = getEnv("KOMGA_API_KEY", "")
	KavitaURL = getEnv("KAVITA_URL", "")
	KavitaAPIKey = getEnv("KAVITA_API_KEY", "")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func parseIntEnv(key string, fallback int) int {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}
