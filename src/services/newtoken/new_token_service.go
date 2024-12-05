package newtokenservice

import (
	"fmt"

	"github.com/gin-gonic/gin"

	dbservice "github.com/ArsenChick/web-service-gin/services/db"
	"github.com/ArsenChick/web-service-gin/utils"
)

type NewTokenService struct {
	dBService *dbservice.DBService
}

var (
	ErrUserNotFound = dbservice.ErrUserNotFound
	ErrBadStructure = fmt.Errorf("wrong request format")
	ErrNotAGuid     = fmt.Errorf("expected valid UUID")
	ErrInternal     = fmt.Errorf("internal server error")
)

type newTokenRequestBody struct {
	Guid string `json:"guid"`
}

func New(dbs *dbservice.DBService) *NewTokenService {
	return &NewTokenService{dBService: dbs}
}

func (s *NewTokenService) PerformNewTokenLogic(ctx *gin.Context) (map[string]string, error) {
	reqBody, err := checkRequestStructure(ctx)
	if err != nil {
		return handleError(ctx, err)
	}

	guid, err := utils.ParseGUIDFromString(reqBody.Guid)
	if err != nil {
		return handleError(ctx, ErrNotAGuid)
	}

	_, err = s.dBService.CheckUserPresentByGUID(guid)
	if err != nil {
		if err == dbservice.ErrUserNotFound {
			return handleError(ctx, ErrUserNotFound)
		} else {
			return handleError(ctx, ErrInternal)
		}
	}

	respBody, rtString := utils.GetNewTokensResponseAndRefreshTokenStr(ctx.ClientIP(), guid)
	rtBcryptHash, err := utils.CreateNewBcryptToken(rtString)
	if err != nil {
		return handleError(ctx, ErrInternal)
	}
	if err = s.dBService.UpdateRefreshTokenHashByGUID(guid, rtBcryptHash); err != nil {
		return handleError(ctx, ErrInternal)
	}
	return respBody, nil
}

func handleError(ctx *gin.Context, err error) (map[string]string, error) {
	if err == ErrInternal {
		// To help me understand what's going on
		fmt.Printf("Error: %s ", err.Error())
	}
	ctx.Abort()
	return nil, err
}

func checkRequestStructure(ctx *gin.Context) (*newTokenRequestBody, error) {
	var requestBody newTokenRequestBody
	err := ctx.BindJSON(&requestBody)
	if err != nil || requestBody.Guid == "" {
		return nil, ErrBadStructure
	}
	return &requestBody, nil
}
