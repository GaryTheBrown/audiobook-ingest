package audiobook

import (
	"audiobook-ingest/data"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (ab *AudioBook) Path(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var m4bFiles []string
	var splitAudioFiles []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(file.Name()))
		filePath := filepath.Join(path, file.Name())

		switch ext {
		case ".m4b":
			m4bFiles = append(m4bFiles, filePath)
		case ".mp3", ".m4a", ".flac", ".aac":
			splitAudioFiles = append(splitAudioFiles, filePath)
		}
	}

	m4bFound := len(m4bFiles) > 0
	splitFound := len(splitAudioFiles) > 0

	switch {
	case !m4bFound && !splitFound:
		return nil

	case m4bFound && !splitFound:
		for _, m4b := range m4bFiles {
			data.Queue(m4b, filepath.Base(m4b), false)
		}
		return nil

	case !m4bFound && splitFound:
		bookGroups := ab.analyzeFileGroupingHeuristics(splitAudioFiles)
		if len(bookGroups) == 1 {
			name := filepath.Base(path)
			data.Queue(path, name, true)
		} else {
			ab.processMultiBookSplits(path, bookGroups)
		}

	case m4bFound && splitFound:
		if err := ab.handleMixedM4BAndSplitTracks(path, m4bFiles, splitAudioFiles); err != nil {
			return err
		}
	}

	return nil
}

func (ab *AudioBook) processMultiBookSplits(basePath string, groups map[string][]string) {
	for bookKey := range groups {
		simulatedNewFolder := filepath.Join(basePath, bookKey)

		data.Queue(simulatedNewFolder, bookKey, true)
	}
}

func (ab *AudioBook) analyzeFileGroupingHeuristics(files []string) map[string][]string {
	groups := make(map[string][]string)
	garbageTrimPattern := regexp.MustCompile(`(?i)(?:part|pt|track|cd|\b)\s*#?\d+\s*$|[\(\[][0-9]+[\)\]]\s*$`)

	for _, file := range files {
		name := filepath.Base(file)
		ext := filepath.Ext(name)
		nameWithoutExt := strings.TrimSuffix(name, ext)

		cleanKey := garbageTrimPattern.ReplaceAllString(nameWithoutExt, "")
		cleanKey = strings.Trim(cleanKey, " -_[]()")
		if cleanKey == "" {
			cleanKey = "Unclassified Track Compilation"
		}

		groups[cleanKey] = append(groups[cleanKey], file)
	}

	if len(groups) > 1 {
		firstKey := ""
		allMatch := true
		for k := range groups {
			if firstKey == "" {
				firstKey = k
				continue
			}
			if !strings.HasPrefix(k, firstKey[:len(firstKey)/2]) {
				allMatch = false
				break
			}
		}

		if allMatch {
			unifiedGroups := make(map[string][]string)
			parentFolder := filepath.Base(filepath.Dir(files[0]))
			for _, file := range files {
				unifiedGroups[parentFolder] = append(unifiedGroups[parentFolder], file)
			}
			return unifiedGroups
		}
	}

	return groups
}

func (ab *AudioBook) handleMixedM4BAndSplitTracks(basePath string, m4bFiles []string, looseTracks []string) error {
	looseGroups := ab.analyzeFileGroupingHeuristics(looseTracks)

	existingM4Bs := make(map[string]string)
	for _, m4b := range m4bFiles {
		m4bName := strings.TrimSuffix(filepath.Base(m4b), filepath.Ext(m4b))
		normalizedM4B := ab.normalizeComparisonString(m4bName)
		existingM4Bs[normalizedM4B] = m4b
	}

	for looseBookKey, tracks := range looseGroups {
		normalizedLooseKey := ab.normalizeComparisonString(looseBookKey)

		foundMatch := ab.findMatchingM4B(normalizedLooseKey, existingM4Bs)

		if foundMatch {

			for _, track := range tracks {
				log.Printf("[SIMULATION LOG] loose track prune: Would permanently delete redundant track -> %s", filepath.Base(track))
				// os.Remove(track)
			}
		} else {
			isolatedSubFolder := filepath.Join(basePath, looseBookKey)
			log.Printf("[SIMULATION LOG] folder create: Would isolate mismatched files into subfolder -> %s", isolatedSubFolder)
			// os.Mkdir(isolatedSubFolder, 0755)
			for _, track := range tracks {
				dest := filepath.Join(isolatedSubFolder, filepath.Base(track))
				log.Printf("[SIMULATION LOG] file Move: Would migrate loose track '%s' to -> %s", filepath.Base(track), dest)
				// os.Rename(track, dest)
			}

			data.Queue(isolatedSubFolder, looseBookKey, true)
		}
	}

	for _, m4b := range m4bFiles {

		data.Queue(m4b, filepath.Base(m4b), false)
	}

	return nil
}

func (ab *AudioBook) findMatchingM4B(looseKey string, existingM4Bs map[string]string) bool {
	if _, exactMatch := existingM4Bs[looseKey]; exactMatch {
		return true
	}

	for m4bKey := range existingM4Bs {
		if strings.Contains(looseKey, m4bKey) || strings.Contains(m4bKey, looseKey) {
			return true
		}
	}
	return false
}

func (ab *AudioBook) normalizeComparisonString(input string) string {
	output := strings.ToLower(input)

	prefixRegex := regexp.MustCompile(`^\d+\s*-\s*|^(?:disc|cd)\s*\d+\s*-_*`)
	output = prefixRegex.ReplaceAllString(output, "")

	output = strings.Trim(output, " -_[]()")
	output = strings.ReplaceAll(output, " ", "")

	if output == "" {
		return "unclassifiedcompilation"
	}
	return output
}
