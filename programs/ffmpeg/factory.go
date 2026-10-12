package ffmpeg

func New() *FFmpeg {
	return &FFmpeg{
		inputs: make([]string, 0),
	}
}

// ffmpeg -f concat -safe 0 -i [manifestFileLocation] -c copy -y [output]
func NewConcatCopy(manifestFileLocation, output string) *FFmpeg {
	runner := New()
	runner.ConcatMode()
	runner.AddInput(manifestFileLocation)
	runner.CopyGlobal()
	runner.Output(output)
	return runner
}

// ffmpeg -f concat -safe 0 -i [manifestFileLocation] -c:a [codec] -b:a [bitrate]k -vn -y [output]
func NewConcatAudioTranscode(manifestFileLocation, output string, codec AudioType, bitrate string) *FFmpeg {
	runner := New()
	runner.ConcatMode()
	runner.AddInput(manifestFileLocation)
	runner.TranscodeAudio(codec, bitrate)
	runner.GlobalNoVideo()
	runner.Output(output)
	return runner
}
