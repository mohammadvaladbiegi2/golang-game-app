package authhttphandler

import (
	"gamegolang/repository/mysql"
	userservice "gamegolang/service/user_service"
	"net/http"

	"github.com/labstack/echo"
)

func Register(c echo.Context) error {

	body := new(userservice.RegisterRequest)

	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}
	mysqlRepo := mysql.NewDB()
	userRepo := userservice.RegisterService{Repo: mysqlRepo}

	respons, RegisterUserError := userRepo.Register(*body)
	if RegisterUserError.HaveError() {

		return c.JSON(RegisterUserError.MetaDataError().StatusCode, RegisterUserError.Jsonmessage())

	}

	return c.JSON(http.StatusOK, echo.Map{
		"ID":          respons.ID,
		"Name":        respons.Name,
		"PhoneNumber": respons.PhoneNumber,
		"Password":    respons.Password,
	})

}
