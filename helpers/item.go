package helpers

import (
	"audiobook-ingest/data"
	"audiobook-ingest/enum"
	"audiobook-ingest/status"
	"audiobook-ingest/structs"
)

var (
	SendMessage func(itemID string, sendType enum.SendType, payloadObj any)
)

func SendBase(i *data.Item) {
	basePayload := i.Base()
	SendMessage(i.ID(), enum.CreateItem, basePayload)
}

func SetStatus(i *data.Item, statusIn status.Ingest) {
	i.SetStatus(statusIn)
	id := i.ID()
	if SendMessage != nil {
		payload := structs.ChangeStatus{
			Status: statusIn,
		}
		SendMessage(id, enum.ChangeStatus, payload)
	}

	switch statusIn {
	case status.Detecting:
		data.DetectingQueue <- id
	case status.Converting:
		data.ConversionQueue <- id
	case status.Verifying:
		data.VerificationQueue <- id
	case status.Verified:
		data.ExportQueue <- id
	default:
	}
}

func RemoveSearchMetadata(i *data.Item) {
	i.RemoveSearchMetadata()
	if SendMessage != nil {
		SendMessage(i.ID(), enum.RemoveMetadata, nil)
	}
}

func RemoveChoices(i *data.Item) {
	i.RemoveChoices()
	if SendMessage != nil {
		SendMessage(i.ID(), enum.RemoveChoices, nil)
	}
}

func RemoveDetactingLogs(i *data.Item) {
	i.RemoveDetactingLogs()

	if SendMessage != nil {
		SendMessage(i.ID(), enum.DetectingRemoveLogs, nil)
	}
}

func SetErrorMessage(i *data.Item, errorMessage error) {
	i.SetErrorMessage(errorMessage)
	if SendMessage != nil {
		failedMessage := structs.FailedMessage{
			Message: errorMessage.Error(),
		}
		SendMessage(i.ID(), enum.FailedMessage, failedMessage)
	}
}

func ClearErrorMessage(i *data.Item) {
	i.ClearErrorMessage()
	if SendMessage != nil {
		failedMessage := structs.FailedMessage{
			Message: "",
		}
		SendMessage(i.ID(), enum.FailedMessage, failedMessage)
	}
}
