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
	if LoginError.HaveError() {
		return c.JSON(LoginError.MetaDataError().StatusCode, LoginError.Jsonmessage())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})

}
