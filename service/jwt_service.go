package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mhmmmdrivaldhi/go-book-api/config"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JwtService interface {
	GenerateToken(userId int, email, role string) (string, error)
	ValidateToken(tokenStr string) (*Claims, error)
}

type jwtService struct {
	cfg config.ApiConfig
}

func (js *jwtService) GenerateToken(userId int, email, role string) (string, error) {
	var app config.AppConfig

	claims := &Claims{
		UserID: userId,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(js.cfg.AccessTokenLifeTime) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer: app.ApplicatonName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedMethod, err := token.SignedString([]byte(js.cfg.JwtSignatureKey))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return signedMethod, nil
}

func (js *jwtService) ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error)  {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(js.cfg.JwtSignatureKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}

		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token not valid yet")
		}
		return nil, fmt.Errorf("token parsing/validation error: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func NewJwtService(cfg config.ApiConfig) JwtService {
	return &jwtService{
		cfg: cfg,
	}
}