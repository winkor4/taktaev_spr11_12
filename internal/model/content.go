package model

import (
	"context"
	"io"
)

type ContentRepo interface {
	AddContent(ctx context.Context, key, dataType string) error
	GetContent(ctx context.Context, name string) (EncContent, error)
	ContentList(ctx context.Context) ([]string, error)
	DeleteContent(ctx context.Context, name string) error
	UpdateContent(ctx context.Context, key, dataType string) error
}

// Интерфейс, который ползволяет прочитать тело запроса в зависимости от типа данных
type AddContentRequest interface {
	JSONDecode(body io.ReadCloser) error
	SchemaToStorageData(user User, dataType string) (StorageData, error)
}
