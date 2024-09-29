package utils

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrorSetupNotCompleted                  Error = "setup not completed"
	ErrorInvalidPath                        Error = "invalid path"
	ErrorFileNotFound                       Error = "file not found"
	ErrorFileProviderNotFound               Error = "file provider not found"
	ErrorElasticSearchScriptNotAcknowledged Error = "elasticsearch script not acknowledged"
	ErrorPermissionDenied                   Error = "permission denied"
)
