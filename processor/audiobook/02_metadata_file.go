package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/status"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func (ab *AudioBook) checkNameParsing(item *data.Item) (foundID bool, foundMeta bool) {
	foundID = false
	foundMeta = false
	tagMap := make(map[string]string)
	seenSegments := make(map[string]bool)

	processTextNode := func(rawText string) []string {
		if ext := filepath.Ext(rawText); ext != "" {
			rawText = strings.TrimSuffix(rawText, ext)
		}

		rawText = strings.ReplaceAll(rawText, " - ", "|")
		rawText = strings.ReplaceAll(rawText, "-", "|")
		rawText = strings.ReplaceAll(rawText, "[", "|")
		rawText = strings.ReplaceAll(rawText, "]", "|")
		rawText = strings.ReplaceAll(rawText, "(", "|")
		rawText = strings.ReplaceAll(rawText, ")", "|")
		rawText = strings.ReplaceAll(rawText, "_", "|")

		rawParts := strings.Split(rawText, "|")
		var cleanParts []string

		noiseRegex := regexp.MustCompile(`(?i)\b(part\s*\d+|pt\s*\d+|cd\s*\d+|vbr|aac|mp3|128kbps|unabridged)\b`)
		spaceCollapse := regexp.MustCompile(`\s+`)

		for _, part := range rawParts {

			cleaned := noiseRegex.ReplaceAllString(part, "")
			cleaned = spaceCollapse.ReplaceAllString(cleaned, " ")
			cleaned = strings.TrimSpace(cleaned)

			if cleaned != "" {
				cleanParts = append(cleanParts, cleaned)
			}
		}
		return cleanParts
	}

	var rawPool []string
	rawPool = append(rawPool, localNameClean(item.Name()))

	currentDir := filepath.Dir(item.FullPathName())
	for i := 0; i < 3; i++ {
		dirName := filepath.Base(currentDir)
		if dirName == "." || dirName == "/" || dirName == "" || strings.Contains(strings.ToLower(dirName), "import") {
			break
		}
		rawPool = append(rawPool, dirName)
		currentDir = filepath.Dir(currentDir)
	}

	segmentIndex := 1
	for _, rawNode := range rawPool {
		fragments := processTextNode(rawNode)
		for _, frag := range fragments {

			normFrag := strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(frag, ""))
			if normFrag == "" {
				continue
			}

			isDuplicate := false
			for processed := range seenSegments {
				if strings.Contains(processed, normFrag) || strings.Contains(normFrag, processed) {
					isDuplicate = true
					break
				}
			}

			if !isDuplicate {
				tagMap[fmt.Sprintf("Segment %d", segmentIndex)] = frag
				seenSegments[normFrag] = true
				segmentIndex++
			}
		}
	}

	if len(tagMap) == 0 {
		tagMap["Folder Path"] = item.Path()
	}

	item.AddMetadata(tagMap)
	item.SetStatus(status.MultiChoice)

	foundMeta = true
	return
}

func localNameClean(s string) string {
	if ext := filepath.Ext(s); ext != "" {
		return strings.TrimSuffix(s, ext)
	}
	return s
}
