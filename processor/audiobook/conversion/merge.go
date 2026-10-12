package conversion

import (
	"audiobook-ingest/programs/ffmpeg"
	"audiobook-ingest/programs/ffprobe"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TrackData struct {
	Path     string
	Duration time.Duration
	Bitrate  uint64
}

var chunkRegexp = regexp.MustCompile(`(\d+|\D+)`)

func NaturalLess(str1, str2 string) bool {
	chunks1 := chunkRegexp.FindAllString(strings.ToLower(str1), -1)
	chunks2 := chunkRegexp.FindAllString(strings.ToLower(str2), -1)

	maxLen := len(chunks1)
	if len(chunks2) < maxLen {
		maxLen = len(chunks2)
	}

	for i := 0; i < maxLen; i++ {
		c1, c2 := chunks1[i], chunks2[i]

		if c1[0] >= '0' && c1[0] <= '9' && c2[0] >= '0' && c2[0] <= '9' {
			n1, err1 := strconv.ParseUint(c1, 10, 64)
			n2, err2 := strconv.ParseUint(c2, 10, 64)

			if err1 == nil && err2 == nil {
				if n1 != n2 {
					return n1 < n2
				}
				continue
			}
		}

		if c1 != c2 {
			return c1 < c2
		}
	}

	return len(chunks1) < len(chunks2)
}

func ReadFolder(bookFolderPath string) ([]TrackData, error) {
	var tracks []TrackData
	probeEngine := ffprobe.New()

	err := filepath.Walk(bookFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".mp3" || ext == ".m4a" {
				probeOutput, probeErr := probeEngine.Run(path)
				if probeErr != nil {
					return nil
				}
				track := TrackData{Path: path}
				if probeOutput.Format != nil {
					if secs, pErr := strconv.ParseFloat(probeOutput.Format.Duration, 64); pErr == nil {
						track.Duration = time.Duration(secs * float64(time.Second))
					}
					if br, pErr := strconv.ParseUint(probeOutput.Format.BitRate, 10, 64); pErr == nil {
						track.Bitrate = br
					}
				}
				tracks = append(tracks, track)
			}
		}
		return nil
	})
	return tracks, err
}

func ExecuteAudioMerge(bookFolderPath, intermediateM4BPath string, onProgress ffmpeg.ProgressCallback) (string, error) {
	tracks, err := ReadFolder(bookFolderPath)
	if err != nil || len(tracks) == 0 {
		return "", fmt.Errorf("zero raw master tracks (.mp3/.m4a) discovered inside directory: %s", bookFolderPath)
	}

	sort.Slice(tracks, func(i, j int) bool {
		return NaturalLess(filepath.Base(tracks[i].Path), filepath.Base(tracks[j].Path))
	})

	allNativeM4A := true
	var totalDuration time.Duration
	var combinedBitrate uint64
	var validProbesCount uint64

	for _, t := range tracks {
		totalDuration += t.Duration
		if t.Bitrate > 0 {
			combinedBitrate += t.Bitrate
			validProbesCount++
		}
		if strings.ToLower(filepath.Ext(t.Path)) != ".m4a" {
			allNativeM4A = false
		}
	}

	manifestPath := filepath.Join(bookFolderPath, "ffmpeg_concat_queue.txt")
	var manifestContent strings.Builder
	for _, t := range tracks {
		abs, err := filepath.Abs(t.Path)
		if err != nil {
			abs = t.Path
		}
		escapedPath := strings.ReplaceAll(abs, "'", "'\\''")
		fmt.Fprintf(&manifestContent, "file '%s'\n", escapedPath)
	}

	if err := os.WriteFile(manifestPath, []byte(manifestContent.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to save text list manifest layout file: %w", err)
	}
	defer os.Remove(manifestPath)

	var cmdEngine *ffmpeg.FFmpeg
	if allNativeM4A {
		cmdEngine = ffmpeg.NewConcatCopy(manifestPath, intermediateM4BPath)
	} else {
		targetBitrate := "128k"
		if validProbesCount > 0 {
			avgKbps := (combinedBitrate / validProbesCount) / 1000
			switch {
			case avgKbps <= 64:
				targetBitrate = "64k"
			case avgKbps <= 96:
				targetBitrate = "96k"
			case avgKbps <= 128:
				targetBitrate = "128k"
			default:
				targetBitrate = "192k"
			}
		}
		cmdEngine = ffmpeg.NewConcatAudioTranscode(manifestPath, intermediateM4BPath, ffmpeg.Aaac, targetBitrate)
	}

	if err := cmdEngine.Run(totalDuration, onProgress); err != nil {
		return "", fmt.Errorf("ffmpeg worker engine execution failed: %w", err)
	}

	return intermediateM4BPath, nil
}
