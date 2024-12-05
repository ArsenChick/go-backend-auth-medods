package refreshtokenservice

import (
	"fmt"

	"github.com/gin-gonic/gin"

	dbservice "github.com/ArsenChick/web-service-gin/services/db"
	"github.com/ArsenChick/web-service-gin/services/mailer"
	"github.com/ArsenChick/web-service-gin/utils"
)

type RefreshTokenService struct {
	dBService     *dbservice.DBService
	mailerService *mailer.MailerService
}

var (
	ErrInvalidGuid    = fmt.Errorf("unauthorized")
	ErrRefreshInvalid = fmt.Errorf("invalid refresh token")
	ErrInternal       = fmt.Errorf("internal server error")
)

func New(dbs *dbservice.DBService) *RefreshTokenService {
	ms := &mailer.MailerService{}
	return &RefreshTokenService{dBService: dbs, mailerService: ms}
}

func (s *RefreshTokenService) PerformRefreshTokenLogic(ctx *gin.Context) (
	map[string]string, error) {

	guidStr := ctx.GetString("guid")
	sentIp := ctx.GetString("iss_ip")
	refreshToken := ctx.GetString("refresh_token")

	guid, err := utils.ParseGUIDFromString(guidStr)
	if err != nil {
		return handleError(ctx, ErrInvalidGuid)
	}

	s.dBService.BeginTransaction()
	defer s.dBService.RollbackTransaction()

	storedHash, err := s.dBService.GetRefreshTokenHashByGUIDTx(guid)
	if err != nil {
		if err == dbservice.ErrUserNotFound {
			return handleError(ctx, ErrInvalidGuid)
		} else {
			return handleError(ctx, ErrInternal)
		}
	}

	isRefreshValid, err := utils.CompareBcryptHash(storedHash, refreshToken)
	if !isRefreshValid {
		return handleError(ctx, ErrRefreshInvalid)
	} else if err != nil {
		return handleError(ctx, ErrInternal)
	}

	respBody, rtString := utils.GetNewTokensResponseAndRefreshTokenStr(ctx.ClientIP(), guid)
	rtBcryptHash, err := utils.CreateNewBcryptToken(rtString)
	if err != nil {
		return handleError(ctx, ErrInternal)
	}
	if err = s.dBService.UpdateRefreshTokenHashByGUIDTx(guid, rtBcryptHash); err != nil {
		return handleError(ctx, ErrInternal)
	}

	if ctx.ClientIP() != sentIp {
		email, err := s.dBService.GetMailByGUIDTx(guid)
		if err == nil {
			s.mailerService.SendWarningEmail(email)
		}
	}

	s.dBService.CommitTransaction()
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
