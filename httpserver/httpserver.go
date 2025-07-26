package httpserver

import (
	"errors"
	authhttphandler "gamegolang/httpserver/httphandler/auth-http-handler"
	userhttphandler "gamegolang/httpserver/httphandler/user-http-handler"
	"log/slog"
	"net/http"

	"github.com/labstack/echo"
)

func Server() {
	e := echo.New()

	// userGroup := e.Group("/users")
	// userGroup.Get("/add",handler)   /users/add
	// userGroup.Get("/remove",handler) /users/remove
	// userGroup.POST("/create",handler) /users/create

	e.GET("/profile", userhttphandler.GetProfile)
	e.POST("/register", authhttphandler.Register)

	if err := e.Start(":5000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}

}
