package httphandler

import (
	"net/http"

	"github.com/labstack/echo"
)

func Hello(c echo.Context) error {

	return c.JSON(http.StatusOK, echo.Map{
		"message": "hello from server OK",
	})
}
