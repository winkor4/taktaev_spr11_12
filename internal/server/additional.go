package server

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Требования регистрации
type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

// Ключ для создания токена
var jwtKey = []byte("secret_key")

// Генерирует и возвращает хэш
func hash(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	return string(bytes), err
}

// Проверяет соответствие пароля и хэша
func checkHash(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

// Создает куку для авторизации
func authToken(login string) (*http.Cookie, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expirationTime,
	}, nil
}

// Возвращает логин авторизованного пользователя
func userFromCtx(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(keyUser).(string)
	if !ok {
		return "", ok
	}
	return user, ok
}
