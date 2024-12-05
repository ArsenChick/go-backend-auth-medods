package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret-key")

func GetNewTokensResponseAndRefreshTokenStr(ip string,
	guid *uuid.UUID) (map[string]string, string) {

	expiryDate := time.Now().Add(time.Hour * 240).Unix()

	jwtTokenClaims := jwt.MapClaims{}
	jwtTokenClaims["iss"] = "arseniy"
	jwtTokenClaims["exp"] = expiryDate
	jwtTokenClaims["ip"] = ip

	jwtTokenClaims["guid"] = guid
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtTokenClaims)
	accessTokenStr, _ := accessToken.SignedString(jwtKey)
	delete(jwtTokenClaims, "guid")

	jwtTokenClaims["jwt-access"] = accessTokenStr
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtTokenClaims)
	refreshTokenStr, _ := refreshToken.SignedString(jwtKey)
	refreshTokenBase64Str := base64.StdEncoding.EncodeToString([]byte(refreshTokenStr))

	response := map[string]string{
		"access":  accessTokenStr,
		"refresh": refreshTokenBase64Str,
	}
	return response, refreshTokenStr
}

func VerifyJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer("arseniy"))

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func CheckTokenPairValidity(accessTokenString string, refreshTokenString string) (jwt.MapClaims, error) {
	accessTokenClaims, err := VerifyJWTToken(accessTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid access token")
	}
	requestTokenClaims, err := VerifyJWTToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}
	if requestTokenClaims["jwt-access"] != accessTokenString {
		return nil, fmt.Errorf("not a valid token pair")
	}
	return accessTokenClaims, nil
}

func CreateNewBcryptToken(refreshToken string) ([]byte, error) {
	sha256HashedToken := sha256.Sum256([]byte(refreshToken))
	bcryptHashedToken, err := bcrypt.GenerateFromPassword(
		sha256HashedToken[:], bcrypt.MinCost)
	if err != nil {
		return nil, err
	} else {
		return bcryptHashedToken, nil
	}
}

func CompareBcryptHash(refreshTokenHash, refreshToken string) (bool, error) {
	sha256HashedToken := sha256.Sum256([]byte(refreshToken))
	err := bcrypt.CompareHashAndPassword([]byte(refreshTokenHash), sha256HashedToken[:])
	if err != nil {
		return false, err
	}
	return true, nil
}

func ParseGUIDFromString(guidStr string) (*uuid.UUID, error) {
	guid, err := uuid.Parse(guidStr)
	if err != nil {
		return nil, fmt.Errorf("not a GUID")
	}
	return &guid, nil
}
