package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/processor/audiobook/structs"
	"audiobook-ingest/programs/ffprobe"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (ab *AudioBook) Verify(item *data.Item) error {
	bookFolder := item.FullPathName()
	finalM4BPath := filepath.Join(bookFolder, fmt.Sprintf("%s.m4b", item.Name()))
	probeEngine := ffprobe.New()

	fileInfo, err := os.Stat(finalM4BPath)
	if err != nil {
		return fmt.Errorf("output file not found on disk: %w", err)
	}

	if fileInfo.Size() < 10240 {
		return fmt.Errorf("compiled file is empty or truncated (%d bytes)", fileInfo.Size())
	}

	probeOutput, err := probeEngine.Run(finalM4BPath)
	if err != nil {
		return fmt.Errorf("ffprobe structural analysis crashed: %w", err)
	}

	hasAudioTrack := false
	var actualFileDuration time.Duration

	for _, stream := range probeOutput.Streams {
		if stream.CodecType == "audio" {
			hasAudioTrack = true
			if stream.Duration != "" {
				if secs, err := strconv.ParseFloat(stream.Duration, 64); err == nil {
					actualFileDuration = time.Duration(secs * float64(time.Second))
				}
			}
			break
		}
	}

	if !hasAudioTrack {
		return fmt.Errorf("container houses zero extractable audio streams")
	}

	if actualFileDuration == 0 && probeOutput.Format != nil && probeOutput.Format.Duration != "" {
		if secs, err := strconv.ParseFloat(probeOutput.Format.Duration, 64); err == nil {
			actualFileDuration = time.Duration(secs * float64(time.Second))
		}
	}

	metadata, ok := item.FinalMetadata().(*structs.AudioBookMetadata)
	if !ok || metadata == nil {
		return fmt.Errorf("verification failed: target item is missing final metadata structures in memory")
	}

	expectedDuration := metadata.ExpectedDuration()

	if expectedDuration > 0 {
		variance := expectedDuration - actualFileDuration
		if variance < 0 {
			variance = -variance
		}

		if variance > (2 * time.Second) {
			return fmt.Errorf("timeline duration mismatch: expected %v, got %v",
				expectedDuration.Round(time.Second), actualFileDuration.Round(time.Second))
		}
	}

	return nil
}
