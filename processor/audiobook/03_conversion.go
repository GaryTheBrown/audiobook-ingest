package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/enum"
	"audiobook-ingest/processor/audiobook/conversion"
	"audiobook-ingest/status"
	"audiobook-ingest/structs"
	"fmt"
	"os"
	"path/filepath"
)

func (ab *AudioBook) Convert(item *data.Item) error {
	bookFolder := item.FullPathName()
	intermediateMergeM4B := filepath.Join(bookFolder, "intermediate_merged_stream.m4b")
	finalM4BPath := filepath.Join(bookFolder, fmt.Sprintf("%s.m4b", item.Name()))

	mergeProgressCallback := func(percent float64) {
		if data.SendMessage != nil {
			progressPayload := structs.ConversionProgressMessage{
				Percent: percent,
				Text:    "Merging audio tracks ...",
			}
			data.SendMessage(item.ID(), enum.ConversionProgress, progressPayload)
		}
	}

	_, err := conversion.ExecuteAudioMerge(bookFolder, intermediateMergeM4B, mergeProgressCallback)
	if err != nil {
		_ = os.Remove(intermediateMergeM4B)
		item.SetStatus(status.Failed)
		return fmt.Errorf("audio merge track compilation failed: %w", err)
	}

	taggingProgressCallback := func(percent float64) {
		if data.SendMessage != nil {
			progressPayload := structs.ConversionProgressMessage{
				Percent: percent,
				Text:    "Embedding metadata tags and chapters...",
			}
			data.SendMessage(item.ID(), enum.ConversionProgress, progressPayload)
		}
	}

	err = conversion.ApplyMetadataAndCoverArt(intermediateMergeM4B, finalM4BPath, item, taggingProgressCallback)

	_ = os.Remove(intermediateMergeM4B)

	if err != nil {
		item.SetStatus(status.Failed)
		return fmt.Errorf("container asset metadata injection failed: %w", err)
	}

	if _, err := os.Stat(finalM4BPath); os.IsNotExist(err) {
		item.SetStatus(status.Failed)
		return fmt.Errorf("compiled target file missing from output path: %s", finalM4BPath)
	}

	if data.SendMessage != nil {
		completionPayload := structs.ConversionProgressMessage{
			Percent: 100.0,
			Text:    "Conversion complete!",
		}
		data.SendMessage(item.ID(), enum.ConversionProgress, completionPayload)
	}

	item.SetStatus(status.Verifying)
	return nil
}
