package server

import (
	"context"
)

// Возвращает логин авторизованного пользователя
func userFromCtx(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(keyUser).(string)
	if !ok {
		return "", ok
	}
	return user, ok
}
