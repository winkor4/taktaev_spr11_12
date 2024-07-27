// Функции для взаимодействия с базой данной
package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
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
func (db *DB) AddUser(ctx context.Context, login string, pass string) (bool, error) {

	result, err := db.db.ExecContext(ctx, queryAddUser,
		login,
		pass)
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
