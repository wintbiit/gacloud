package utils

type Error string

func (e Error) Error() string {
	return string(e)
}

var ErrorFileNotFound Error = "file not found"
var ErrorFileProviderNotFound Error = "file provider not found"
