package userhttphandler

import (
	"gamegolang/pkg/jwt"
	"gamegolang/repository/mysql"
	userservice "gamegolang/service/user_service"
	"net/http"

	"github.com/labstack/echo"
)

func GetProfile(c echo.Context) error {

	AuthorizationToken := c.Request().Header.Get("Authorization")

	if AuthorizationToken == "" {
		return c.String(http.StatusUnauthorized, "You do not have access")
	}

	stringAuthorization := AuthorizationToken[7:]
	VerifyResult, vError := jwt.VerifyToken(stringAuthorization)
	if vError.HaveError() {
		return c.JSON(vError.MetaDataError().StatusCode, vError.Jsonmessage())
	}

	mysqlRepo := mysql.NewDB()
	LoginRepo := userservice.LoginService{
		Repo: mysqlRepo,
	}

	userName, profileError := LoginRepo.GetProfile(VerifyResult.ID)
	if profileError.HaveError() {
		return c.JSON(profileError.MetaDataError().StatusCode, profileError.Jsonmessage())
	}

	return c.JSON(http.StatusOK, echo.Map{"Name": userName.Name})

}
