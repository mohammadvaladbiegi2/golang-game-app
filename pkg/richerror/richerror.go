package richerror

import (
	"github.com/labstack/echo"
)

type RichError struct {
	wrappedError error
	statusCode   int
	message      string
	metaData     map[string]interface{}
}

func NewError(errorData ...interface{}) RichError {
	r := RichError{}

	for _, value := range errorData {
		switch value.(type) {
		case error:
			r.wrappedError = value.(error)
		case int:
			r.statusCode = value.(int)
		case string:
			r.message = value.(string)
		case map[string]interface{}:
			r.metaData = value.(map[string]interface{})

		}
	}

	return r
}

func (r *RichError) MetaDataError() struct {
	Message    string
	StatusCode int
} {

	metaDataToClient := struct {
		Message    string
		StatusCode int
	}{
		Message:    r.message,
		StatusCode: r.statusCode,
	}

	return metaDataToClient
}

func (r *RichError) Jsonmessage() echo.Map {
	return echo.Map{
		"message":    r.message,
		"statusCode": r.statusCode,
	}
}

func (r RichError) Error() string {
	return r.message
}

func (r RichError) HaveError() bool {
	if r.statusCode == 0 {
		return false
	}

	return true
}
