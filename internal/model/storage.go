package model

import "context"

type StorageRepo interface {
	Ping(ctx context.Context) error
	AddUser(ctx context.Context, data User) (bool, error)
	GetPass(ctx context.Context, login string) (string, error)
	GetKey(ctx context.Context, login string) (string, error)
	AddContent(ctx context.Context, sData StorageData) error
	GetContent(ctx context.Context, name, user string) (EncContent, error)
	ContentList(ctx context.Context, user string) ([]string, error)
	DeleteContent(ctx context.Context, name, user string) error
	UpdateContent(ctx context.Context, sData StorageData) error
}
