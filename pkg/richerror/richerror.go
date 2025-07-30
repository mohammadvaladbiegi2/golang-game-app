package richerror

type RichError struct {
	WrappedError error
	StatusCode   uint16
	Message      string
	MetaData     map[string]interface{}
}

func NewError(errorData RichError) RichError {
	return RichError{
		WrappedError: errorData.WrappedError,
		StatusCode:   errorData.StatusCode,
		Message:      errorData.Message,
		MetaData:     errorData.MetaData,
	}
}

func (r RichError) Error() string {
	return r.Message
}
