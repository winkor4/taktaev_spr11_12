package users

import (
	"context"

	"github.com/winkor4/taktaev_spr11_12/internal/crypto"
	"github.com/winkor4/taktaev_spr11_12/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserManager struct {
	login    string
	password string
	key      string
	db       model.StorageRepo
}

func NewUserManager(l, p, k string, db model.StorageRepo) *UserManager {
	return &UserManager{
		login:    l,
		password: p,
		key:      k,
		db:       db,
	}
}

func (u *UserManager) GetLogin() string {
	return u.login
}

func (u *UserManager) AddUser(ctx context.Context) (bool, error) {
	hash, err := hash(u.password)
	if err != nil {
		return false, err
	}

	encryptionSK := crypto.RandStr(16)
	encryptionSK, err = crypto.Encrypt(encryptionSK, u.key)
	if err != nil {
		return false, err
	}

	conflict, err := u.db.AddUser(ctx, model.GerUserModel(u.login, hash, encryptionSK))
	if err != nil {
		return false, err
	}

	return conflict, nil
}

func (u *UserManager) CheckAuth(ctx context.Context) (bool, error) {
	hash, err := u.GetPass(ctx)
	if err != nil {
		return false, err
	}
	return checkHash(u.password, hash), nil
}

func (u *UserManager) GetKey(ctx context.Context) (string, error) {
	encKey, err := u.db.GetKey(ctx, u.login)
	if err != nil {
		return "", err
	}
	return encKey, nil
}

func (u *UserManager) GetPass(ctx context.Context) (string, error) {
	hash, err := u.db.GetPass(ctx, u.login)
	if err != nil {
		return "", err
	}
	return hash, nil
}

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
