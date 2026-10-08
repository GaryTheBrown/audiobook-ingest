package audiobook

import (
	"audiobook-ingest/data"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
)

func (ab *AudioBook) checkID3Tags(item *data.Item) (foundID bool, foundMeta bool) {
	foundID = false
	foundMeta = false
	fullPath := item.FullPathName()
	targetFile := fullPath
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		err = filepath.Walk(fullPath, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !info.IsDir() {
				ext := strings.ToLower(filepath.Ext(path))
				if ext == ".mp3" || ext == ".m4a" || ext == ".m4b" || ext == ".flac" {
					targetFile = path
					return filepath.SkipDir
				}
			}
			return nil
		})
		if err != nil || targetFile == fullPath {
			return
		}
	}

	f, err := os.Open(targetFile)
	if err != nil {
		item.AddDetactingLog("⚠️", fmt.Sprintf("Failed to open audio track for tag parsing: %s", filepath.Base(targetFile)))
		f.Close()
		return
	}

	m, err := tag.ReadFrom(f)
	f.Close()
	if err != nil {
		item.AddDetactingLog("ℹ", "No embedded ID3/MP4 metadata header tags located inside audio tracks.")
		return
	}

	item.AddDetactingLog("🎵", fmt.Sprintf("Inspecting embedded embedded tags inside track: %s", filepath.Base(targetFile)))

	tagMap := make(map[string]string)

	seenValues := make(map[string]bool)

	addCoreTag := func(label, value string) {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			tagMap[label] = trimmed
			seenValues[strings.ToLower(trimmed)] = true
		}
	}

	addCoreTag("Title", m.Title())
	addCoreTag("Author", m.Artist())
	addCoreTag("Album", m.Album())
	addCoreTag("Genre", m.Genre())
	addCoreTag("Comment", m.Comment())
	if m.Year() > 0 {
		addCoreTag("Year", strconv.Itoa(m.Year()))
	}

	if rawData := m.Raw(); rawData != nil {
		for key, interfaceValue := range rawData {
			if interfaceValue == nil {
				continue
			}

			var cleanValue string
			switch v := interfaceValue.(type) {
			case string:
				cleanValue = strings.TrimSpace(v)
			case []byte:
				cleanValue = strings.TrimSpace(string(v))
			default:
				cleanValue = strings.TrimSpace(fmt.Sprintf("%v", v))
			}

			if cleanValue == "" {
				continue
			}

			if seenValues[strings.ToLower(cleanValue)] {
				continue
			}

			normKey := strings.ToLower(key)
			switch {
			case normKey == "tcom" || normKey == "composer":
				tagMap["Composer"] = cleanValue
				seenValues[strings.ToLower(cleanValue)] = true
			case normKey == "text" || strings.Contains(normKey, "narrator"):
				tagMap["Narrator"] = cleanValue
				seenValues[strings.ToLower(cleanValue)] = true
			case normKey == "tpub" || normKey == "publisher":
				tagMap["Publisher"] = cleanValue
				seenValues[strings.ToLower(cleanValue)] = true
			case normKey == "tit3" || strings.Contains(normKey, "subtitle"):
				tagMap["Subtitle"] = cleanValue
				seenValues[strings.ToLower(cleanValue)] = true
			default:
				if !strings.HasPrefix(key, "unknown") && len(tagMap) < 30 {
					tagMap["Raw_"+key] = cleanValue
					seenValues[strings.ToLower(cleanValue)] = true
				}
			}
		}
	}

	item.AddMetadata(tagMap)

	asinRegex := regexp.MustCompile(`\bB[A-Z0-9]{9}\b`)
	for _, key := range []string{"Comment", "Title", "Album"} {
		if field, exists := tagMap[key]; exists {
			if matched := asinRegex.FindString(strings.ToUpper(field)); matched != "" {
				item.AddDetactingLog("🎯", fmt.Sprintf("Audible ASIN match isolated inside embedded audio tag '%s': %s", key, matched))
				item.SetTokenID(matched)
				foundID = true
				return
			}
		}
	}

	if tagMap["Title"] != "" {
		item.AddDetactingLog("✔", fmt.Sprintf("Extracted embedded tags -> Title: '%s' | Artist: '%s'", tagMap["Title"], tagMap["Author"]))
	}
	foundMeta = true
	return
}
