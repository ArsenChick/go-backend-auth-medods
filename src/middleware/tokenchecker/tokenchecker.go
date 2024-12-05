package tokenchecker

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ArsenChick/web-service-gin/utils"
)

func TokenCheckerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenString := c.GetHeader("Access-Token")
		if accessTokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing access token"})
			c.Abort()
			return
		}

		refreshTokenStringBase64 := c.GetHeader("Refresh-Token")
		if refreshTokenStringBase64 == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
			c.Abort()
			return
		}

		refreshTokenString, err := base64.StdEncoding.DecodeString(refreshTokenStringBase64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "malformed refresh token"})
			c.Abort()
			return
		}

		claims, err := utils.CheckTokenPairValidity(accessTokenString, string(refreshTokenString))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("guid", claims["guid"])
		c.Set("iss_ip", claims["ip"])
		c.Set("refresh_token", string(refreshTokenString))
		c.Next()
	}
}
