package audiobook

import (
	"audiobook-ingest/data"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var asinRegex = regexp.MustCompile(`\bB[A-Z0-9]{9}\b`)

func (ab *AudioBook) checkMetadataJson(item *data.Item) (foundID bool, foundMeta bool) {
	foundID = false
	foundMeta = false
	searchDir := item.FullPathName()
	fileInfo, err := os.Stat(searchDir)
	if err == nil && !fileInfo.IsDir() {
		searchDir = filepath.Dir(searchDir)
	}

	var jsonPath string
	err = filepath.Walk(searchDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".json") {
			lowerName := strings.ToLower(info.Name())
			if strings.Contains(lowerName, "metadata") || strings.Contains(lowerName, "manifest") {
				jsonPath = path
				return filepath.SkipDir
			}
		}
		return nil
	})

	if err != nil || jsonPath == "" {
		return
	}

	fileBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		item.AddDetactingLog("⚠️", fmt.Sprintf("Failed to read discovered JSON file: %s", filepath.Base(jsonPath)))
		return
	}

	var rawMetadata map[string]any
	if err := json.Unmarshal(fileBytes, &rawMetadata); err != nil {
		item.AddDetactingLog("❌", fmt.Sprintf("JSON syntax parsing failed for: %s", filepath.Base(jsonPath)))
		return
	}

	item.AddDetactingLog("📝", fmt.Sprintf("Successfully parsed and flattened local metadata manifest: %s", filepath.Base(jsonPath)))

	flattenedResult := make(map[string]string)
	matchedASIN := flattenAndScanForASIN(rawMetadata, flattenedResult, "")
	item.AddMetadata(flattenedResult)

	tagMap := SortJson(rawMetadata)
	if matchedASIN == "" {
		for _, fieldVal := range tagMap {
			if match := asinRegex.FindString(strings.ToUpper(fieldVal)); match != "" {
				matchedASIN = match
				break
			}
		}
	}

	if matchedASIN != "" {
		item.AddDetactingLog("🎯", fmt.Sprintf("Audible ASIN match isolated inside JSON text: %s", matchedASIN))
		foundID = item.SetTokenID(matchedASIN)
		return
	}

	if tagMap["Title"] != "" && tagMap["Author"] != "" {
		item.AddDetactingLog("✔", fmt.Sprintf("Extracted base fields from JSON -> Title: '%s' | Author: '%s'", tagMap["Title"], tagMap["Author"]))
		foundMeta = true
		return
	}

	item.AddDetactingLog("ℹ", "Local JSON file scanned, but lacked an explicit ASIN identifier or core tag fields.")
	return
}

func flattenAndScanForASIN(obj map[string]any, result map[string]string, parentKey string) string {
	var foundASIN string

	for key, val := range obj {
		deepPathKey := key
		if parentKey != "" {
			deepPathKey = parentKey + "." + key
		}

		if val == nil {
			result[deepPathKey] = "—"
			continue
		}

		switch v := val.(type) {
		case map[string]any:
			asinFromNested := flattenAndScanForASIN(v, result, deepPathKey)
			if asinFromNested != "" && foundASIN == "" {
				foundASIN = asinFromNested
			}

		case []any:
			var lines []string
			for _, item := range v {
				if item == nil {
					continue
				}
				if innerMap, ok := item.(map[string]any); ok {
					bytes, _ := json.Marshal(innerMap)
					lines = append(lines, string(bytes))
				} else {
					lines = append(lines, strings.TrimSpace(fmt.Sprintf("%v", item)))
				}
			}
			if len(lines) > 0 {
				joinedVal := strings.Join(lines, ", ")
				result[deepPathKey] = joinedVal

				if match := asinRegex.FindString(strings.ToUpper(joinedVal)); match != "" && foundASIN == "" {
					foundASIN = match
				}
			}

		case string:
			trimmedStr := strings.TrimSpace(v)
			result[deepPathKey] = trimmedStr

			if match := asinRegex.FindString(strings.ToUpper(trimmedStr)); match != "" && foundASIN == "" {
				foundASIN = match
			}

		case float64:
			result[deepPathKey] = strconv.FormatFloat(v, 'f', -1, 64)

		case bool:
			result[deepPathKey] = strconv.FormatBool(v)

		default:
			trimmedStr := strings.TrimSpace(fmt.Sprintf("%v", v))
			result[deepPathKey] = trimmedStr

			if match := asinRegex.FindString(strings.ToUpper(trimmedStr)); match != "" && foundASIN == "" {
				foundASIN = match
			}
		}
	}

	return foundASIN
}

func SortJson(payload map[string]any) map[string]string {

	tagMap := make(map[string]string)

	extractJSONField := func(standardLabel string, jsonKeys ...string) string {
		for _, key := range jsonKeys {
			if val, exists := payload[key]; exists && val != nil {
				var strVal string
				switch v := val.(type) {
				case string:
					strVal = strings.TrimSpace(v)
				case float64:
					strVal = strconv.FormatFloat(v, 'f', 0, 64)

				case []any:
					var parts []string
					for _, item := range v {
						if item != nil {
							parts = append(parts, strings.TrimSpace(fmt.Sprintf("%v", item)))
						}
					}
					if len(parts) > 0 {
						strVal = strings.Join(parts, ", ")
					}

				default:
					strVal = strings.TrimSpace(fmt.Sprintf("%v", v))
				}
				if strVal != "" {
					tagMap[standardLabel] = strVal
					return strVal
				}
			}

			for k, val := range payload {
				if strings.EqualFold(k, key) && val != nil {
					var strVal string
					switch v := val.(type) {
					case string:
						strVal = strings.TrimSpace(v)
					case float64:
						strVal = strconv.FormatFloat(v, 'f', 0, 64)
					case []any:
						var parts []string
						for _, item := range v {
							if item != nil {
								parts = append(parts, strings.TrimSpace(fmt.Sprintf("%v", item)))
							}
						}
						if len(parts) > 0 {
							strVal = strings.Join(parts, ", ")
						}
					default:
						strVal = strings.TrimSpace(fmt.Sprintf("%v", val))
					}
					if strVal != "" {
						tagMap[standardLabel] = strVal
						return strVal
					}
				}
			}
		}
		return ""
	}

	extractJSONField("Title", "title", "bookTitle", "book_title")
	extractJSONField("Author", "author", "artist", "authorName", "author_name")
	extractJSONField("Album", "album", "series_name", "series")
	extractJSONField("Genre", "genre")
	extractJSONField("Narrator", "narrator", "narrator_name", "narratorName")
	extractJSONField("Publisher", "publisher", "publisher_name")
	extractJSONField("Year", "year", "release_date", "pub_date", "publishDate")
	extractJSONField("Comment", "comment", "description", "summary")

	toTitleCase := func(s string) string {
		s = strings.ReplaceAll(strings.ReplaceAll(s, "_", " "), "-", " ")
		words := strings.Fields(s)
		for i, word := range words {
			if len(word) > 0 {
				words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			}
		}
		return strings.Join(words, " ")
	}

	for k, val := range payload {
		displayKey := toTitleCase(k)

		if _, exists := tagMap[displayKey]; !exists && val != nil {
			var strVal string
			switch v := val.(type) {
			case string:
				strVal = strings.TrimSpace(v)
			case bool:
				strVal = strconv.FormatBool(v)
			case float64:
				strVal = strconv.FormatFloat(v, 'f', -1, 64)

			case []any:
				var lines []string
				for _, item := range v {
					if item == nil {
						continue
					}

					switch inner := item.(type) {
					case map[string]any:
						title := ""
						start := ""
						end := ""

						for subKey, subVal := range inner {
							lowerSubKey := strings.ToLower(subKey)
							valStr := strings.TrimSpace(fmt.Sprintf("%v", subVal))
							if strings.Contains(valStr, ".") {
								if f, err := strconv.ParseFloat(valStr, 64); err == nil {
									valStr = strconv.FormatFloat(f, 'f', 2, 64)
								}
							}

							switch lowerSubKey {
							case "title", "name":
								title = valStr
							case "start":
								start = valStr
							case "end":
								end = valStr
							}
						}

						if title != "" {
							lines = append(lines, fmt.Sprintf("%s (Start: %ss, End: %ss)", title, start, end))
						} else {
							lines = append(lines, fmt.Sprintf("Start: %ss, End: %ss", start, end))
						}

					default:
						lines = append(lines, strings.TrimSpace(fmt.Sprintf("%v", item)))
					}
				}

				if len(lines) > 0 {
					if strings.ToLower(k) == "chapters" {
						strVal = strings.Join(lines, "<br>")
					} else {
						strVal = strings.Join(lines, ", ")
					}
				}

			default:
				strVal = strings.TrimSpace(fmt.Sprintf("%v", v))
			}

			if strVal != "" && len(strVal) < 5000 {
				tagMap[displayKey] = strVal
			}
		}
	}
	return tagMap
}
