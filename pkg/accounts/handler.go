package accounts

import (
	"database/sql"
	"errors"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterHandler(e *echo.Echo, accountService Service) {
	e.GET("/:aspspId/accounts/cid/:cid", callAccounts(accountService))
	e.GET("/:aspspId/accounts/:accountId/cid/:cid", callAccounts(accountService))
}

func callAccounts(s Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		rid := c.Response().Header().Get(echo.HeaderXRequestID)
		aspspId := c.Param("aspspId")
		if aspspId == "" {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "aspspId can't be empty"))
		}

		cid := c.Param("cid")
		if aspspId == "" {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "cid can't be empty"))
		}

		accountId := c.Param("accountId")
		var res string
		var err error
		if accountId == "" {
			res, err = s.Accounts(cid, aspspId)
		} else {
			res, err = s.Account(cid, aspspId, accountId)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, err.Error()))
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, err.Error()))
		}

		return c.JSON(http.StatusOK, res)
	}
}
