package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dbservice "github.com/ArsenChick/web-service-gin/services/db"
	newtokenservice "github.com/ArsenChick/web-service-gin/services/newtoken"
	refreshtokenservice "github.com/ArsenChick/web-service-gin/services/refreshtoken"
)

type Controller struct {
	newTokenService     *newtokenservice.NewTokenService
	refreshTokenService *refreshtokenservice.RefreshTokenService
}

func New(dbs *dbservice.DBService) *Controller {
	nts := newtokenservice.New(dbs)
	rts := refreshtokenservice.New(dbs)
	return &Controller{
		newTokenService:     nts,
		refreshTokenService: rts,
	}
}

func (c *Controller) HandleNewTokenRequest(ctx *gin.Context) {
	var (
		statusCode int
		response   any
	)

	jsonObject, err := c.newTokenService.PerformNewTokenLogic(ctx)
	if err != nil {
		response = gin.H{"message": err.Error()}
	} else {
		response = jsonObject
	}

	switch err {
	case newtokenservice.ErrUserNotFound:
		statusCode = http.StatusBadRequest
	case newtokenservice.ErrBadStructure:
		statusCode = http.StatusBadRequest
	case newtokenservice.ErrNotAGuid:
		statusCode = http.StatusUnprocessableEntity
	case newtokenservice.ErrInternal:
		statusCode = http.StatusInternalServerError
	default:
		statusCode = http.StatusOK
	}

	ctx.JSON(statusCode, response)
}

func (c *Controller) HandleRefreshTokenRequest(ctx *gin.Context) {
	var (
		statusCode int
		response   any
	)

	jsonObject, err := c.refreshTokenService.PerformRefreshTokenLogic(ctx)
	if err != nil {
		response = gin.H{"message": err.Error()}
	} else {
		response = jsonObject
	}

	switch err {
	case refreshtokenservice.ErrInvalidGuid:
		statusCode = http.StatusUnauthorized
	case refreshtokenservice.ErrRefreshInvalid:
		statusCode = http.StatusUnauthorized
	case refreshtokenservice.ErrInternal:
		statusCode = http.StatusInternalServerError
	default:
		statusCode = http.StatusOK
	}

	ctx.JSON(statusCode, response)
}
