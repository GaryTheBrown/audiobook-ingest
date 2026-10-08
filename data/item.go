package data

import (
	"audiobook-ingest/enum"
	"audiobook-ingest/status"
	"audiobook-ingest/structs"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	SendMessage  func(itemID string, sendType enum.SendType, payloadObj any)
	CheckTokenID func(tokenID string) (exists bool, name string, imgSrc string, metadata map[string]string)
)

type Item struct {
	mu                 sync.RWMutex
	id                 string
	addedAt            int64
	fullPathName       string
	name               string
	path               string
	image              string
	requiredConversion bool
	status             status.Ingest

	tokenID       string
	foundImgSrc   string
	foundMetadata map[string]string
	metadata      map[string]string
	detectingLog  []structs.DetectingLogs
	choices       []structs.Choice
}

func NewItem(fullPath, name, path string, reqConversion bool) *Item {
	return &Item{
		id:                 fmt.Sprintf("%d", time.Now().UnixNano()),
		addedAt:            time.Now().Unix(),
		fullPathName:       fullPath,
		name:               name,
		path:               path,
		requiredConversion: reqConversion,
		status:             status.Detecting,
		metadata:           map[string]string{},
		detectingLog:       []structs.DetectingLogs{},
		choices:            []structs.Choice{},
	}
}

func (i *Item) SendBase() {
	i.mu.RLock()
	basePayload := structs.CreateItem{
		AddedAt: i.addedAt,
		Name:    i.name,
		Path:    i.path,
		Image:   i.image,
		Status:  i.status,
	}
	i.mu.RUnlock()
	SendMessage(i.id, enum.CreateItem, basePayload)
}

func (i *Item) ID() string {
	return i.id
}

func (i *Item) AddedAt() int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.addedAt
}

func (i *Item) FullPathName() string {
	return i.fullPathName
}

func (i *Item) Name() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.name
}

func (i *Item) Path() string {
	return i.path
}

func (i *Item) Image() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if i.image != "" {
		return i.image

	}
	return "/favicon.svg"
}

func (i *Item) SetImage(image string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.image = image
}

func (i *Item) RequiredConversion() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.requiredConversion
}

func (i *Item) SetRequiredConversion(req bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.requiredConversion = req
}

func (i *Item) Status() status.Ingest {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.status
}

func (i *Item) SetStatus(statusIn status.Ingest) {
	i.mu.Lock()
	i.status = statusIn
	i.mu.Unlock()
	payload := structs.ChangeStatus{
		Status: statusIn,
	}
	if SendMessage != nil {
		SendMessage(i.id, enum.ChangeStatus, payload)
	}

	switch statusIn {
	case status.Detecting:
		DetectingQueue <- i.id
	case status.Converting:
		ConversionQueue <- i.id
	case status.Verified:
		VerificationQueue <- i.id
	default:
	}
}

func (i *Item) TokenID() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.tokenID
}
func (i *Item) setTokenID(id string, isUserOverride bool) {
	trimmedID := strings.TrimSpace(id)

	exists, name, imgSrc, data := CheckTokenID(trimmedID)
	if !exists {
		if isUserOverride {
			SendMessage(i.id, enum.AcceptToken, structs.AcceptToken{Accept: false})
		}
		return
	}

	i.mu.Lock()
	i.tokenID = trimmedID
	i.name = name
	i.foundImgSrc = imgSrc
	i.foundMetadata = data
	i.mu.Unlock()

	if SendMessage != nil && trimmedID != "" {
		acceptPayload := structs.AcceptToken{
			Accept:  true,
			TokenID: trimmedID,
			Name:    name,
			ImgSrc:  imgSrc,
			Data:    data,
		}

		SendMessage(i.id, enum.AcceptToken, acceptPayload)
	}
}
func (i *Item) SetTokenIDManual(id string) {
	i.setTokenID(id, true)
}

func (i *Item) SetTokenID(id string) {
	i.setTokenID(id, false)
}

func (i *Item) FoundImgSrc() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.foundImgSrc
}
func (i *Item) FoundMetadata() map[string]string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.foundMetadata
}

func (i *Item) AddMetadata(metadata map[string]string) {
	i.mu.Lock()
	i.metadata = metadata
	i.mu.Unlock()
	if len(metadata) > 0 {
		payload := structs.AddMetadata{
			Fields: metadata,
		}
		if SendMessage != nil {
			SendMessage(i.id, enum.AddMetadata, payload)
		}
	}
}

func (i *Item) RemoveMetadata() {
	i.mu.Lock()
	i.metadata = map[string]string{}
	i.mu.Unlock()

	if SendMessage != nil {
		SendMessage(i.id, enum.RemoveMetadata, nil)
	}
}

func (i *Item) Metadata() map[string]string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.metadata
}

func (i *Item) SetChoices(choices []structs.Choice) {
	i.mu.Lock()
	i.choices = choices
	i.mu.Unlock()
	choicesPayload := structs.AddChoices{
		Choices: choices,
	}
	if SendMessage != nil {
		SendMessage(i.id, enum.AddChoices, choicesPayload)
	}
}

func (i *Item) RemoveChoices() {
	i.mu.Lock()
	i.choices = []structs.Choice{}
	i.mu.Unlock()

	if SendMessage != nil {
		SendMessage(i.id, enum.RemoveChoices, nil)
	}
}

func (i *Item) Choices() []structs.Choice {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.choices
}

func (i *Item) AddDetactingLog(icon, value string) {
	payload := structs.DetectingLogs{
		Icon:  icon,
		Value: value,
	}

	i.mu.Lock()
	i.detectingLog = append(i.detectingLog, payload)
	i.mu.Unlock()

	if SendMessage != nil {
		SendMessage(i.id, enum.DetectingAddLog, payload)
	}
}

func (i *Item) RemoveDetactingLogs() {
	i.mu.Lock()
	i.detectingLog = []structs.DetectingLogs{}
	i.mu.Unlock()

	if SendMessage != nil {
		SendMessage(i.id, enum.DetectingRemoveLogs, nil)
	}
}

func (i *Item) DetectingLogs() []structs.DetectingLogs {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.detectingLog
}
