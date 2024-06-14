package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/model/entity"
	"time"
)

func Parse(tokenString string, cfg config.JWTConfig) (*entity.JwtData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err

	}

	id := int(claims["id"].(float64))
	userRole := claims["user_role"].(string)
	expiration, err := time.Parse(time.RFC3339, claims["expiration"].(string))
	if err != nil {
		return nil, err
	}
	res := &entity.JwtData{
		ID:         id,
		UserRole:   userRole,
		Expiration: expiration,
	}
	return res, nil
}

func Generate(userId int, userRole string, cfg config.JWTConfig) (string, error) {

	expiration := time.Second * cfg.Expiration
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         userId,
		"user_role":  userRole,
		"expiration": time.Now().Add(expiration),
	})

	tokenString, err := token.SignedString([]byte(cfg.Secret))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
