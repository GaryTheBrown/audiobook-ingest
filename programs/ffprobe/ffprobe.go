package ffprobe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
)

type FFprobeFormat struct {
	Duration string            `json:"duration"`
	BitRate  string            `json:"bit_rate"`
	Size     string            `json:"size"`
	Tags     map[string]string `json:"tags,omitempty"`
}

type FFprobeStream struct {
	Index      int    `json:"index"`
	CodecType  string `json:"codec_type"`
	CodecName  string `json:"codec_name"`
	BitRate    string `json:"bit_rate,omitempty"`
	Duration   string `json:"duration,omitempty"`
	SampleRate string `json:"sample_rate,omitempty"`
	Channels   int    `json:"channels,omitempty"`
}

type FFprobeOutput struct {
	Format  *FFprobeFormat  `json:"format,omitempty"`
	Streams []FFprobeStream `json:"streams,omitempty"`
}

type FFprobe struct {
	binaryPath string
	baseArgs   []string
}

func New() *FFprobe {
	return &FFprobe{
		binaryPath: "ffprobe",
		baseArgs:   []string{"-loglevel", "error", "-show_format", "-show_streams", "-output_format", "json"},
	}
}

func (p *FFprobe) WithCustomBinary(customPath string) *FFprobe {
	if customPath != "" {
		p.binaryPath = customPath
	}
	return p
}

func (p *FFprobe) Run(filePath string) (*FFprobeOutput, error) {
	cmdArgs := make([]string, len(p.baseArgs)+1)
	copy(cmdArgs, p.baseArgs)
	cmdArgs[len(p.baseArgs)] = filePath

	cmd := exec.Command(p.binaryPath, cmdArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffprobe failed for %s: %w | stderr: %s", filepath.Base(filePath), err, stderr.String())
	}

	var output FFprobeOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return nil, fmt.Errorf("failed to map ffprobe json payload for %s: %w", filepath.Base(filePath), err)
	}

	return &output, nil
}
