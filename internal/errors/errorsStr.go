package serviceErrors

// HTTP Status Codes
const (
	StatusNotFound            = 404
	StatusBadRequest          = 400
	StatusInternalServerError = 500
)

// 500
const (
	ErrUploadingUserCode    = "ERROR_WHILE_UPLOADING_USER_CODE"
	ErrCreateDirectoryError = "ERROR_CREATE_DIRECTORY"
	ErrCodeBuild            = "ERROR_CODE_BUILD"
)
