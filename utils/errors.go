package utils

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrorSetupNotCompleted    Error = "setup not completed"
	ErrorFileNotFound         Error = "file not found"
	ErrorFileProviderNotFound Error = "file provider not found"
)
