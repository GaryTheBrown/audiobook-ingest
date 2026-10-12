package ffmpeg

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

func (f *FFmpeg) Run(totalDuration time.Duration, onProgress ProgressCallback) error {
	cmd := exec.Command(binaryPath, f.args...)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to open stderr console tracking pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ffmpeg sub-process failed to initialize: %w", err)
	}

	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() {
		line := scanner.Text()

		if matches := progressRegex.FindStringSubmatch(line); len(matches) > 4 && totalDuration > 0 {
			hours, _ := strconv.Atoi(matches[1])
			mins, _ := strconv.Atoi(matches[2])
			secs, _ := strconv.Atoi(matches[3])

			elapsed := time.Duration(hours)*time.Hour + time.Duration(mins)*time.Minute + time.Duration(secs)*time.Second

			percent := (float64(elapsed) / float64(totalDuration)) * 100.0
			if percent > 100.0 {
				percent = 100.0
			}

			if onProgress != nil {
				onProgress(percent)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		// Log softly, non-fatal to the actual conversion loop process execution
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg process exited with an error code: %w", err)
	}

	return nil
}
