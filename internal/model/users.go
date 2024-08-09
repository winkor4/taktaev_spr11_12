package model

import "context"

type UsersRepo interface {
	GetLogin() string
	GerUserModel() User
	AddUser(ctx context.Context) (bool, error)
	CheckAuth(ctx context.Context) (bool, error)
	GetPass(ctx context.Context) (string, error)
	GetKey(ctx context.Context) (string, error)
}
