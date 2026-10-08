package enum

type SendType int

const (
	CreateItem SendType = iota
	RemoveItem
	AcceptToken
	AddMetadata
	RemoveMetadata
	ChangeStatus
	DetectingAddLog
	DetectingRemoveLogs
	AddChoices
	RemoveChoices
)

type RecievedType int

const (
	Ping RecievedType = iota
	DetectingStuck
	SubmitToken
	FailedRetry
	FailedDelete
	SwitchToManual
	SelectChoice
	ManualSearch
	Rescan
)
