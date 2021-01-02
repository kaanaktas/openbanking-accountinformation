package callback

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterHandler(e *echo.Echo, callbackService Service) {
	e.GET("/callback", processCallBack(callbackService))
}

func processCallBack(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.QueryParam("code")
		state := c.QueryParam("state")

		var httpStatusCode int
		var message string
		err := service.ProcessCallBack(code, state)
		if errors.Is(err, sql.ErrNoRows) {
			httpStatusCode = http.StatusNotFound
		} else {
			switch err {
			case sql.ErrNoRows:
			case nil:
				httpStatusCode = http.StatusOK
				message = "Consent has been authorized successfully."
			default:
				httpStatusCode = http.StatusInternalServerError
				message = "An unexpected error has occurred. " + err.Error()
			}
		}

		return c.JSON(httpStatusCode, message)
	}
}
