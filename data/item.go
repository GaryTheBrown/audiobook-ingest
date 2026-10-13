package data

import (
	"audiobook-ingest/status"
	"audiobook-ingest/structs"
	"fmt"
	"sync"
	"time"
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

	searchMetadata map[string]string
	detectingLog   []structs.DetectingLogs
	choices        []MetadataInterface

	tokenID         string
	finalCoverImage string
	finalMetadata   MetadataInterface

	errorMessage error
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
	}
}

func (i *Item) Base() *structs.CreateItem {
	i.mu.RLock()
	base := structs.CreateItem{
		AddedAt: i.addedAt,
		Name:    i.name,
		Path:    i.path,
		Image:   i.image,
		Status:  i.status,
	}
	i.mu.RUnlock()
	return &base
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

func (i *Item) SetName(name string) {
	i.mu.Lock()
	i.name = name
	i.mu.Unlock()
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
	i.image = image
	i.mu.Unlock()
}

func (i *Item) RequiredConversion() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.requiredConversion
}

func (i *Item) SetRequiredConversion(req bool) {
	i.mu.Lock()
	i.requiredConversion = req
	i.mu.Unlock()
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
}

func (i *Item) TokenID() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.tokenID
}
func (i *Item) SetTokenID(tokenID string) {
	i.mu.Lock()
	i.tokenID = tokenID
	i.mu.Unlock()
}

func (i *Item) AddSearchMetadata(metadata map[string]string) {
	i.mu.Lock()
	i.searchMetadata = metadata
	i.mu.Unlock()

}

func (i *Item) RemoveSearchMetadata() {
	i.mu.Lock()
	i.searchMetadata = nil
	i.mu.Unlock()
}

func (i *Item) Metadata() map[string]string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.searchMetadata
}

func (i *Item) SetChoices(choices []MetadataInterface) {
	i.mu.Lock()
	i.choices = choices
	i.mu.Unlock()
}

func (i *Item) RemoveChoices() {
	i.mu.Lock()
	i.choices = []MetadataInterface{}
	i.mu.Unlock()
}

func (i *Item) Choices() []MetadataInterface {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.choices
}

func (i *Item) ChoicesFront() []structs.Choice {
	i.mu.RLock()
	defer i.mu.RUnlock()
	metadataChoices := i.choices
	var structChoices []structs.Choice
	for _, choice := range metadataChoices {
		structChoices = append(structChoices, choice.ToChoice())
	}
	return structChoices
}

func (i *Item) AddDetactingLog(log structs.DetectingLogs) {
	i.mu.Lock()
	i.detectingLog = append(i.detectingLog, log)
	i.mu.Unlock()
}

func (i *Item) RemoveDetactingLogs() {
	i.mu.Lock()
	i.detectingLog = []structs.DetectingLogs{}
	i.mu.Unlock()
}

func (i *Item) DetectingLogs() []structs.DetectingLogs {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.detectingLog
}

func (i *Item) SetFinalCoverImage(imgsrc string) {
	i.mu.Lock()
	i.finalCoverImage = imgsrc
	i.mu.Unlock()
}
func (i *Item) FinalCoverImage() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.finalCoverImage
}

func (i *Item) FinalMetadata() MetadataInterface {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.finalMetadata
}

func (i *Item) SetFinalMetadata(data MetadataInterface) {
	i.mu.Lock()
	i.finalMetadata = data
	i.mu.Unlock()
}

func (i *Item) SetErrorMessage(errorMessage error) {
	i.mu.Lock()
	i.errorMessage = errorMessage
	i.mu.Unlock()
}

func (i *Item) ErrorMessage() error {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.errorMessage
}

func (i *Item) ClearErrorMessage() {
	i.mu.RLock()
	i.errorMessage = nil
	i.mu.RUnlock()
}
