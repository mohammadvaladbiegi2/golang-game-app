package authhttphandler

import (
	"gamegolang/repository/mysql"
	userservice "gamegolang/service/user_service"
	"net/http"

	"github.com/labstack/echo"
)

func Login(c echo.Context) error {

	var bodyData userservice.LoginCredentials
	if err := c.Bind(&bodyData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid request payload",
		})
	}

	mysqlRepo := mysql.NewDB()
	LoginRepo := userservice.LoginService{
		Repo: mysqlRepo,
	}

	token, LoginError := LoginRepo.Login(bodyData)
	if LoginError != nil {

		switch LoginError.Error() {
		case "user not found":
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "user not found",
			})

		case "password or phone number does not match":
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "password or phone number does not match",
			})

		case "server Error":
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "server Error",
			})

		default:
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": LoginError.Error(),
			})
		}

	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})

}
