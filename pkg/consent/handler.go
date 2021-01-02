package consent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
)

func RegisterHandler(e *echo.Echo, sessionService session.Service, service ServiceRead, facadeService Facade) {
	e.GET("/:aspspId/internal/consent/active", retrieveActiveConsent(sessionService, service))
	e.POST("/:aspspId/account-access-consents/reference/:trackingId", createConsent(sessionService, facadeService))
	e.GET("/:aspspId/account-access-consents/:cid", getConsent(facadeService))
}

func retrieveActiveConsent(sessionService session.Service, service ServiceRead) echo.HandlerFunc {
	return func(c echo.Context) error {
		rid := c.Response().Header().Get(echo.HeaderXRequestID)
		aspspId := c.Param("aspspId")
		if aspspId == "" {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "aspspId can't be empty"))
		}

		bearerToken, err := extractBearerToken(c.Request().Header.Get(api.Authorization))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, err.Error()))
		}

		sessionResp, err := sessionService.FindByInternalAccessToken(bearerToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, err.Error()))
		}

		if sessionResp == nil {
			return c.JSON(http.StatusNotFound, api.JsonResponse(rid, "couldn't find the session"))
		}

		consents, err := service.FindAuthorisedConsentByUserIdAndTppId(sessionResp.UserId, sessionResp.TppId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, err.Error()))
		}

		if consents == nil {
			return c.JSON(http.StatusNotFound, api.JsonResponse(rid, "couldn't find any consent"))
		}

		jsonPayload, err := json.Marshal(consents)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, err.Error()))
		}

		return c.JSON(http.StatusOK, string(jsonPayload))
	}
}

func createConsent(sessionService session.Service, proxy Facade) echo.HandlerFunc {
	return func(c echo.Context) error {
		rid := c.Response().Header().Get(echo.HeaderXRequestID)
		aspspId := c.Param("aspspId")
		if aspspId == "" {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "aspspId can't be empty"))
		}
		trackingId := c.Param("trackingId")
		if trackingId == "" {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "reference can't be empty"))
		}

		bearerToken, err := extractBearerToken(c.Request().Header.Get(api.Authorization))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, err.Error()))
		}

		sessionData, err := sessionService.FindByInternalAccessToken(bearerToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, err.Error()))
		}

		consent := &ObReadConsent{}
		if err := c.Bind(consent); err != nil {
			log.Error(err)
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "invalid request. couldn't retrieve the consent details"))
		}

		res, err := proxy.CreateConsent(sessionData.ReferenceId, trackingId, aspspId, consent)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, err.Error()))
		}

		return c.JSON(http.StatusOK, res)
	}
}

func getConsent(proxy Facade) echo.HandlerFunc {
	return func(c echo.Context) error {
		rid := c.Response().Header().Get(echo.HeaderXRequestID)
		aspspId := c.Param("aspspId")
		if aspspId == "" {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "aspspId can't be empty"))
		}

		cid := c.Param("cid")
		if cid == "" {
			return c.JSON(http.StatusBadRequest, api.JsonResponse(rid, "cid can't be empty"))
		}

		res, err := proxy.GetConsent(cid, aspspId)
		switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusNotFound, api.JsonResponse(rid, "couldn't find the consent"))
		case nil:
			return c.JSON(http.StatusOK, res)
		default:
			return c.JSON(http.StatusInternalServerError, api.JsonResponse(rid, err.Error()))
		}
	}
}

func extractBearerToken(authorizationHeader string) (string, error) {
	if strings.HasPrefix(authorizationHeader, "Bearer") {
		return authorizationHeader[7:], nil
	} else {
		return "", fmt.Errorf("unexpected error while extracting the bearerToken")
	}
}
