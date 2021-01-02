package session

import (
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterHandler(e *echo.Echo, service Service) {
	e.GET("/internal/access/initiate/:userId/:tppId/:tid", initiateSession(service))
}

func initiateSession(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		rid := c.Response().Header().Get(echo.HeaderXRequestID)
		userId := c.Param("userId")
		if userId == "" {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "userId can't be empty"))
		}

		tppId := c.Param("tppId")
		if tppId == "" {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "tppId can't be empty"))
		}

		tid := c.Param("tid")
		if tid == "" {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "tid can't be empty"))
		}

		response, err := service.InitiateSession(userId, tppId, tid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, err.Error()))
		}

		return c.JSON(http.StatusOK, response)
	}
}
