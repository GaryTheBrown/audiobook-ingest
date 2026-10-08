package audiobook

import (
	"audiobook-ingest/data"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	durationRegex = regexp.MustCompile(`(?i)total target audiobook duration:\s*(\d{2}):(\d{2}):(\d{2})`)
	ffmpegRegex   = regexp.MustCompile(`size=.*time=(\d{2}):(\d{2}):(\d{2})\.\d{2}.*speed=\s*([\d.]+x)`)
	task1Regex    = regexp.MustCompile(`(?i)(extracting chapter markers|chapters)`)
	task2Regex    = regexp.MustCompile(`(?i)(embedding cover art|optimising mp4 atoms|atoms structure)`)
	task3Regex    = regexp.MustCompile(`(?i)(metadata tag parameters locked|clean-up complete)`)
)

type ConversionState struct {
	mu           sync.Mutex
	totalSeconds int
}

func (ab *AudioBook) CompileConversionCmd(item *data.Item) (*exec.Cmd, error) {
	inputPath := item.FullPathName()
	outputM4B := filepath.Join(filepath.Dir(inputPath), item.Name()+".m4b")

	var audioFiles []string
	err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".mp3" || ext == ".m4a" || ext == ".flac" || ext == ".aac" {
				audioFiles = append(audioFiles, path)
			}
		}
		return nil
	})
	if err != nil || len(audioFiles) == 0 {
		return nil, fmt.Errorf("zero operational audio tracks located for conversion processing paths")
	}

	args := append([]string{}, audioFiles...)
	args = append(args, fmt.Sprintf("--output=%s", outputM4B))
	if item.TokenID() != "" {
		args = append(args, fmt.Sprintf("--metadata-id=%s", item.TokenID()))
	}

	return exec.Command("m4b-merge", args...), nil
}

func (ab *AudioBook) ParseProgressLine(line string, rawState any) (int, string, bool) {
	lineLower := strings.ToLower(line)

	state, _ := rawState.(*ConversionState)
	if state == nil {
		return -1, "", false
	}

	if dMatches := durationRegex.FindStringSubmatch(line); len(dMatches) > 3 {
		hours, _ := strconv.Atoi(dMatches[1])
		mins, _ := strconv.Atoi(dMatches[2])
		secs, _ := strconv.Atoi(dMatches[3])

		state.mu.Lock()
		state.totalSeconds = (hours * 3600) + (mins * 60) + secs
		state.mu.Unlock()

		return 0, "● Initializing audio merge streams...", true
	}

	if strings.Contains(lineLower, "size=") && strings.Contains(lineLower, "time=") && strings.Contains(lineLower, "speed=") {
		if fMatches := ffmpegRegex.FindStringSubmatch(line); len(fMatches) > 4 {
			curHours, _ := strconv.Atoi(fMatches[1])
			curMins, _ := strconv.Atoi(fMatches[2])
			curSecs, _ := strconv.Atoi(fMatches[3])
			currentSeconds := (curHours * 3600) + (curMins * 60) + curSecs

			speedStr := fMatches[4]

			state.mu.Lock()
			totalSecs := state.totalSeconds
			state.mu.Unlock()

			pct := 0
			etaStr := "Calculating..."
			if totalSecs > 0 {
				pct = int((float64(currentSeconds) / float64(totalSecs)) * 100)
				if pct > 100 {
					pct = 100
				}

				remainingSecs := totalSecs - currentSeconds
				speedNum, err := strconv.ParseFloat(strings.TrimSuffix(speedStr, "x"), 64)
				if err == nil && speedNum > 0 {
					realSecondsLeft := int(float64(remainingSecs) / speedNum)
					etaHours := realSecondsLeft / 3600
					etaMins := (realSecondsLeft % 3600) / 60
					etaSecs := realSecondsLeft % 60
					etaStr = fmt.Sprintf("%02d:%02d:%02d", etaHours, etaMins, etaSecs)
				}
			}

			metaLabel := fmt.Sprintf("Speed: %s | ETA: %s", speedStr, etaStr)
			return pct, metaLabel, false
		}
		return -1, "", false
	}

	if strings.Contains(lineLower, "audio compilation track finished") {
		return -1, "", false
	}

	if task1Regex.MatchString(line) {
		return 33, "[Task 1/3] Extracting Chapter Alignment Markers...", true
	}
	if task2Regex.MatchString(line) {
		return 66, "[Task 2/3] Compiling Metadata Header Tags & Cover Art...", true
	}
	if task3Regex.MatchString(line) {
		return 100, "[Task 3/3] Finalizing Media Optimizations...", true
	}

	if strings.Contains(lineLower, "[info]") || strings.Contains(lineLower, "error") {
		return -1, "", true
	}

	return -1, "", false
}
