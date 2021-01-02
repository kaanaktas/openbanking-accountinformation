package config

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/security"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"strings"
)

func NewEchoEngine() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// Middlewares
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(validate())

	return e
}

var permittedUri = []string{"/internal", "/callback", "/favicon.ico"}

//Custom middleware to validate requests JWT
func validate() echo.MiddlewareFunc {
	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//check if request uri is in permission list
			rid := c.Response().Header().Get(echo.HeaderXRequestID)
			reqUri := c.Request().RequestURI
			for _, v := range permittedUri {
				if strings.HasPrefix(reqUri, v) {
					return handler(c)
				}
			}

			log.Debug(reqUri + " is not permitted. Checking for JWT validation...")
			keyData, err := security.GetSecretKey(api.InternalSignKey, os.Getenv("INTERNAL_SIGN_KEY"))
			if err != nil {
				log.Error(err)
				return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, "internal server error"))
			}

			var bearerToken string
			authorizationHeader := c.Request().Header.Get(api.Authorization)
			if strings.HasPrefix(authorizationHeader, "Bearer") {
				bearerToken = authorizationHeader[7:]
			}

			if err := security.VerifyJwt(bearerToken, jwt.SigningMethodHS256, keyData); err != nil {
				log.Error(err)
				return c.JSON(http.StatusUnauthorized, api.JsonResponse(rid, "Unauthorized request"))
			}

			return handler(c)
		}
	}
}
