package api

type ErrUploadError struct {
	Err        error
	UploadType string
}

func (e ErrUploadError) Error() string {
	return "failed to upload " + e.UploadType
}
