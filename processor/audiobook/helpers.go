package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/enum"
	"audiobook-ingest/structs"
)

func addDetactingLog(i *data.Item, icon, value string) {
	payload := structs.DetectingLogs{
		Icon:  icon,
		Value: value,
	}

	i.AddDetactingLog(payload)

	if SendMessage != nil {
		SendMessage(i.ID(), enum.DetectingAddLog, payload)
	}
}

func setChoices(i *data.Item, choices []data.MetadataInterface) {
	i.SetChoices(choices)
	choicesPayload := structs.AddChoices{
		Choices: i.ChoicesFront(),
	}
	if SendMessage != nil {
		SendMessage(i.ID(), enum.AddChoices, choicesPayload)
	}
}

func addSearchMetadata(i *data.Item, metadata map[string]string) {
	i.AddSearchMetadata(metadata)

	if len(metadata) > 0 {
		payload := structs.AddMetadata{
			Fields: metadata,
		}
		if SendMessage != nil {
			SendMessage(i.ID(), enum.AddMetadata, payload)
		}
	}
}
