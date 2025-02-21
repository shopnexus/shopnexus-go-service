package util

import (
	"errors"
	"shopnexus-go-service/config"
	"shopnexus-go-service/internal/model"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID int64, role model.Role) (string, error) {
	tokenDuration := time.Duration(config.GetConfig().App.AccessTokenDuration)

	claims := model.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "safe-trade",
			Subject:   strconv.Itoa(int(userID)),
			Audience:  []string{"safe-trade"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.GetConfig().SensitiveKeys.JWTSecret

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateAccessToken(tokenStr string) (claims model.Claims, err error) {
	keys := config.GetConfig().SensitiveKeys
	secret := keys.JWTSecret

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return claims, errors.New("unverified token or token is expired")
	}

	return claims, nil
}
