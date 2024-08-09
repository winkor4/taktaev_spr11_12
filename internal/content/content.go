package content

import (
	"context"

	"github.com/winkor4/taktaev_spr11_12/internal/crypto"
	"github.com/winkor4/taktaev_spr11_12/internal/model"
	"github.com/winkor4/taktaev_spr11_12/internal/users"
)

type ContentManager struct {
	db                 model.StorageRepo
	contentCredentials model.AddContentRequest
	user               string
	key                string
	dataType           string
	name               string
}

type Config struct {
	DB                 model.StorageRepo
	ContentCredentials model.AddContentRequest
	User               string
	Key                string
	DataType           string
	Name               string
}

func NewContentManager(cfg Config) *ContentManager {
	return &ContentManager{
		db:                 cfg.DB,
		contentCredentials: cfg.ContentCredentials,
		user:               cfg.User,
		key:                cfg.Key,
		dataType:           cfg.DataType,
		name:               cfg.Name,
	}
}

func (c *ContentManager) AddContent(ctx context.Context) error {
	sData, err := c.storageData(ctx)
	if err != nil {
		return err
	}
	err = c.db.AddContent(ctx, sData)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContentManager) GetContent(ctx context.Context) (model.EncContent, error) {
	encData, err := c.db.GetContent(ctx, c.name, c.user)
	if err != nil {

		return model.EncContent{}, err
	}
	return encData, nil
}

func (c *ContentManager) ContentList(ctx context.Context) ([]string, error) {
	dataList, err := c.db.ContentList(ctx, c.user)
	if err != nil {
		return nil, err
	}
	return dataList, err
}

func (c *ContentManager) DeleteContent(ctx context.Context) error {
	err := c.db.DeleteContent(ctx, c.name, c.user)
	if err != nil {

		return err
	}
	return nil
}

func (c *ContentManager) UpdateContent(ctx context.Context) error {
	sData, err := c.storageData(ctx)
	if err != nil {
		return err
	}
	err = c.db.UpdateContent(ctx, sData)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContentManager) storageData(ctx context.Context) (model.StorageData, error) {
	userManager := users.NewUserManager(c.user, "", c.key, c.db)

	encKey, err := userManager.GetKey(ctx)
	if err != nil {
		return model.StorageData{}, err
	}

	key, err := crypto.Decrypt(encKey, c.key)
	if err != nil {
		return model.StorageData{}, err
	}

	sData, err := c.contentCredentials.SchemaToStorageData(model.GerUserModel(c.user, "", key), c.dataType)
	if err != nil {
		return model.StorageData{}, err
	}

	return sData, nil
}
