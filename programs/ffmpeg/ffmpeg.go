package ffmpeg

import (
	"fmt"
	"regexp"
)

const binaryPath string = "ffmpeg"

type ProgressCallback func(percent float64)

var progressRegex = regexp.MustCompile(`time=(\d+):(\d+):(\d+)\.(\d+)`)

type FFmpegAudio struct {
	codec   AudioType
	bitrate string
}

type FFmpeg struct {
	inputs []string
	args   []string
}

func (f *FFmpeg) ConcatMode() {
	f.args = append(f.args, "-f", "concat", "-safe", "0")
}

func (f *FFmpeg) AddInput(path string) int {
	f.inputs = append(f.inputs, path)
	f.args = append(f.args, "-i", path)
	return len(f.inputs) - 1
}

func (f *FFmpeg) Output(path string) {
	f.args = append(f.args, "-y", path)
}

func (f *FFmpeg) CopyGlobal() {
	f.args = append(f.args, "-c", "copy")
}

func (f *FFmpeg) CopyStreamType(st StreamType) {
	f.args = append(
		f.args,
		fmt.Sprintf("-c:%s", st.ToVal()),
		"copy",
	)
}

func (f *FFmpeg) CopyStream(st StreamType, streamID int) {
	f.args = append(
		f.args,
		fmt.Sprintf("-c:%s:%d", st.ToVal(), streamID),
		"copy",
	)
}

func (f *FFmpeg) TranscodeAudio(codec AudioType, bitrate string) {
	if codec == ATNIL || bitrate == "" {
		return
	}
	f.args = append(f.args, "-c:a", codec.ToString(), "-b:a", bitrate)
}

func (f *FFmpeg) GlobalNoVideo() {
	f.args = append(f.args, "-vn")
}

func (f *FFmpeg) Map(stream int, st StreamType) {
	f.args = append(f.args, "-map", fmt.Sprintf("%d:%s", stream, st.ToVal()))
}

func (f *FFmpeg) MapMetadata(FileID int) {
	f.args = append(f.args, "-map_metadata", fmt.Sprintf("%d", FileID))
}

func (f *FFmpeg) TranscodeVideo(vt VideoType) {
	f.args = append(f.args, "-c:v", vt.ToString())
}

func (f *FFmpeg) DispositionGlobal(st StreamType, dt DispositionType) {
	f.args = append(
		f.args,
		fmt.Sprintf("-disposition:%s", st.ToVal()),
		dt.ToString(),
	)
}

func (f *FFmpeg) DispositionStream(st StreamType, streamID int, dt DispositionType) {
	f.args = append(
		f.args,
		fmt.Sprintf("-disposition:%s:%d", st.ToVal(), streamID),
		dt.ToString(),
	)
}
