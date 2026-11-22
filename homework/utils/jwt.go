package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/config"
)

var jwtSecret = []byte("my-super-secret-key") // 在生产环境中应该从环境变量读取

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	var expireTime time.Duration
	expire := cfg.JWT.Expire
	if expire != "" {
		expireTime, err = time.ParseDuration(expire)
		if err != nil {
			return "", err
		}
	}
	if expireTime == 0 {
		expireTime = 24 * time.Hour
	}
	jwtSecret = []byte(cfg.JWT.Secret)
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		cfg, err := config.LoadConfig()
		if err != nil {
			return nil, err
		}
		jwtSecret = []byte(cfg.JWT.Secret)
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
