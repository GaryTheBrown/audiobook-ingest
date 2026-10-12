package conversion

import (
	"audiobook-ingest/data"
	"audiobook-ingest/processor/audiobook/structs"
	"audiobook-ingest/programs/ffmpeg"
	"audiobook-ingest/programs/webgetter"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ApplyMetadataAndCoverArt(intermediateM4B, finalM4BPath string, item *data.Item, onProgress ffmpeg.ProgressCallback) error {
	parentDir := filepath.Dir(intermediateM4B)
	metadataTxtPath, err := CreateMetadataFile(parentDir, item)
	defer os.Remove(metadataTxtPath)
	if err != nil {
		return err
	}

	cmdEngine := ffmpeg.New()
	audioFileID := cmdEngine.AddInput(intermediateM4B)
	metadataFileID := cmdEngine.AddInput(metadataTxtPath)
	imageFileID := -1

	coverImg := item.FinalCoverImage()
	if coverImg != "" {
		localCover, err := webgetter.WebGetter(coverImg, parentDir, "cover")
		if err == nil {
			imageFileID = cmdEngine.AddInput(localCover)
		}
	}

	cmdEngine.Map(audioFileID, ffmpeg.Audio)
	cmdEngine.CopyStreamType(ffmpeg.Audio)

	if imageFileID >= 0 {
		cmdEngine.Map(imageFileID, ffmpeg.Video)
		cmdEngine.TranscodeVideo(ffmpeg.Vpng)
		cmdEngine.DispositionGlobal(ffmpeg.Video, ffmpeg.Dattached_pic)
	}

	cmdEngine.MapMetadata(metadataFileID)
	cmdEngine.Output(finalM4BPath)

	if err := cmdEngine.Run(0, onProgress); err != nil {
		return fmt.Errorf("failed to inject metadata assets into container headers: %w", err)
	}

	return nil
}

func escapeMetadataTag(val string) string {
	res := strings.ReplaceAll(val, "=", "\\=")
	res = strings.ReplaceAll(res, ";", "\\;")
	res = strings.ReplaceAll(res, "#", "\\#")
	return strings.ReplaceAll(res, "\n", "\\\n")
}

func CreateMetadataFile(parentDir string, item *data.Item) (string, error) {
	metadata := item.FinalMetadata().(*structs.AudioBookMetadata)
	metadataTxtPath := filepath.Join(parentDir, "ffmpeg_metadata_tags.txt")
	var sb strings.Builder
	sb.WriteString(";FFMETADATA1\n")

	fmt.Fprintf(&sb, "ASIN=%s\n", escapeMetadataTag(item.TokenID()))
	fmt.Fprintf(&sb, "title=%s\n", escapeMetadataTag(item.Name()))
	fmt.Fprintf(&sb, "album=%s\n", escapeMetadataTag(item.Name()))
	fmt.Fprintf(&sb, "comment=%s\n", escapeMetadataTag(metadata.Description))
	fmt.Fprintf(&sb, "publisher=%s\n", escapeMetadataTag(metadata.Publisher))

	if metadata.Subtitle != "" {
		fmt.Fprintf(&sb, "subtitle=%s\n", escapeMetadataTag(metadata.Subtitle))
	}
	if len(metadata.Authors) > 0 {
		fmt.Fprintf(&sb, "artist=%s\n", escapeMetadataTag(strings.Join(metadata.Authors, ", ")))
	}
	if len(metadata.Genres) > 0 {
		fmt.Fprintf(&sb, "genre=%s\n", escapeMetadataTag(strings.Join(metadata.Genres, ", ")))
	}
	if metadata.Year > 0 {
		fmt.Fprintf(&sb, "date=%d\n", metadata.Year)
	}

	var groupingParts []string
	if len(metadata.Narrators) > 0 {
		groupingParts = append(groupingParts, fmt.Sprintf("Narrator: %s", strings.Join(metadata.Narrators, ", ")))
	}
	if metadata.SeriesName != "" {
		seriesInfo := metadata.SeriesName
		if metadata.SeriesPosition != "" {
			seriesInfo = fmt.Sprintf("%s #%s", metadata.SeriesName, metadata.SeriesPosition)
		}
		groupingParts = append(groupingParts, fmt.Sprintf("Series: %s", seriesInfo))
	}
	if len(groupingParts) > 0 {
		fmt.Fprintf(&sb, "grouping=%s\n", escapeMetadataTag(strings.Join(groupingParts, " | ")))
	}

	for idx, ch := range metadata.Chapters {
		sb.WriteString("[CHAPTER]\n")
		sb.WriteString("TIMEBASE=1/1000\n")
		fmt.Fprintf(&sb, "START=%d\n", ch.StartOffsetMs)

		endOffset := ch.StartOffsetMs + ch.LengthMs
		fmt.Fprintf(&sb, "END=%d\n", endOffset)

		titleStr := ch.Title
		if strings.TrimSpace(titleStr) == "" {
			titleStr = fmt.Sprintf("Chapter %d", idx+1)
		}

		fmt.Fprintf(&sb, "title=%s\n", escapeMetadataTag(titleStr))
	}

	if err := os.WriteFile(metadataTxtPath, []byte(sb.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to save intermediate metadata overlay file: %w", err)
	}
	return metadataTxtPath, nil
}
