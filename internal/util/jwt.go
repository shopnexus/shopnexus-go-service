package util

import (
	"errors"
	"fmt"
	"shopnexus-go-service/config"
	"shopnexus-go-service/internal/model"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(account model.Account) (string, error) {
	tokenDuration := time.Duration(config.GetConfig().App.AccessTokenDuration * int64(time.Second))

	claims := model.Claims{
		UserID: account.GetBase().ID,
		Role:   account.GetBase().Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "shopnexus",
			Subject:   strconv.Itoa(int(account.GetBase().ID)),
			Audience:  []string{"shopnexus"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.GetConfig().SensitiveKeys.JWTSecret

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	// test if token is valid
	_, err = ValidateAccessToken(signedToken)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %w", err)
	}

	return signedToken, nil
}

func ValidateAccessToken(tokenStr string) (claims model.Claims, err error) {
	secret := config.GetConfig().SensitiveKeys.JWTSecret

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return claims, err
	}

	if !token.Valid {
		return claims, errors.New("invalid token or token expired")
	}

	return claims, nil
}
