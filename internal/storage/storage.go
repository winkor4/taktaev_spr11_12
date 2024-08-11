// Функции для взаимодействия с базой данной
package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/winkor4/taktaev_spr11_12/internal/model"
)

// DB - база данных
type DB struct {
	db *sql.DB
}

// New - возвращает соединение с базой данных
func New(ctx context.Context, dsn string) (*DB, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	out := new(DB)
	out.db = db

	err = out.Ping(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, migration := range migrations {
		if _, err := tx.Exec(migration); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &DB{db: db}, nil
}

// Ping - проверяет соединение с базой данных
func (db *DB) Ping(ctx context.Context) error {
	if err := db.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// Создает нового пользователя
func (db *DB) AddUser(ctx context.Context, data model.User) (bool, error) {

	result, err := db.db.ExecContext(ctx, queryAddUser,
		data.Login,
		data.Password,
		data.Key)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 0, err
}

// Поиск хэша пароля пользователя
func (db *DB) GetPass(ctx context.Context, login string) (string, error) {

	row := db.db.QueryRowContext(ctx, queryPassword, login)

	pass := new(string)
	err := row.Scan(pass)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return *pass, nil
}

// Поиск encryptionSK пользователя
func (db *DB) GetKey(ctx context.Context, login string) (string, error) {

	row := db.db.QueryRowContext(ctx, queryKey, login)

	pass := new(string)
	err := row.Scan(pass)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return *pass, nil
}

// Запись данных в БД
func (db *DB) AddContent(ctx context.Context, sData model.StorageData) error {

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(queryInsertContent,
		sData.ID,
		sData.User.Login,
		sData.Name,
		sData.ContentType,
		sData.Data,
		sData.DataSK); err != nil {
		err = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		err = tx.Rollback()
		return err
	}

	return nil
}

// Возвращает данные из БД по имени
func (db *DB) GetContent(ctx context.Context, name, user string) (model.EncContent, error) {
	var result model.EncContent
	result.Name = name
	row := db.db.QueryRowContext(ctx, queryGetContent, user, name)

	err := row.Scan(&result.Data, &result.DataSK, &result.ContentType, &result.EncSK)
	if err == sql.ErrNoRows {
		return result, nil
	}
	if err != nil {
		return result, err
	}

	return result, nil
}

// Возвращает список данных пользователя
func (db *DB) ContentList(ctx context.Context, user string) ([]string, error) {
	rows, err := db.db.QueryContext(ctx, queryContentList, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]string, 0)
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		result = append(result, name)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

// Удаляет данные из БД
func (db *DB) DeleteContent(ctx context.Context, name, user string) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(queryDeleteContent, user, name); err != nil {
		err = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		err = tx.Rollback()
		return err
	}

	return nil
}

// Обновление данных в БД
func (db *DB) UpdateContent(ctx context.Context, sData model.StorageData) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(queryUpdateContent,
		sData.ContentType,
		sData.Data,
		sData.DataSK,
		sData.User.Login,
		sData.Name); err != nil {
		err = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		err = tx.Rollback()
		return err
	}

	return nil
}
