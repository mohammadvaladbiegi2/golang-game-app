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
	if RegisterUserError != nil {

		switch RegisterUserError.Error() {
		case "phone number is not valid":
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "phone number is not valid",
			})

		case "phone number is not unique":
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "phone number is not unique",
			})

		case "password length should be greater than 8":
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "password length should be greater than 8",
			})

		case "name length should be greater than 3":
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "name length should be greater than 3",
			})

		default:
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": RegisterUserError.Error(),
			})
		}

	}

	return c.JSON(http.StatusOK, echo.Map{
		"ID":          respons.ID,
		"Name":        respons.Name,
		"PhoneNumber": respons.PhoneNumber,
		"Password":    respons.Password,
	})

}
