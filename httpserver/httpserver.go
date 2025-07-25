package httpserver

import (
	"errors"
	"gamegolang/httpserver/httphandler"
	"log/slog"
	"net/http"

	"github.com/labstack/echo"
)

func Server() {
	e := echo.New()

	// Routes
	e.GET("/", httphandler.Hello)

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}

}
