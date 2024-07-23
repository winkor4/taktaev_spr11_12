// Функции для взаимодействия с базой данной
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
func (db *DB) Register(ctx context.Context, login string, pass string) (bool, error) {

	result, err := db.db.ExecContext(ctx, queryRegister,
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
	if err != nil {
		return "", err
	}

	return *pass, nil
}

// Запись произвольных текстовых данных в БД
func (db *DB) UploadTextData(ctx context.Context, user string, data []model.TextDataList) ([]model.UploadTextDataResponse, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}

	result, err := db.checkTextData(ctx, tx, user, data)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, textData := range result.resultData {
		_, err := tx.ExecContext(ctx, queryUploadTextData,
			textData.ID,
			user,
			textData.Data)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return result.response, nil
}

// Исключает занятые идентификаторы
func (db *DB) checkTextData(ctx context.Context, tx *sql.Tx, user string, data []model.TextDataList) (checkTextDataResult, error) {

	var result checkTextDataResult

	var param string
	for i, textData := range data {
		param = param + fmt.Sprintf("'%s'", textData.ID)
		if i < len(data)-1 {
			param = param + ", "
		}
	}
	query := strings.ReplaceAll(queryConflictID, "$2", param)

	rows, err := tx.QueryContext(ctx, query, user)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	conflictID := make([]string, 0)
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return result, err
		}
		conflictID = append(conflictID, id)
	}

	resultData := make([]model.TextDataList, len(conflictID))
	response := make([]model.UploadTextDataResponse, len(data))
	for _, textData := range data {
		check := true
		var res model.UploadTextDataResponse
		for _, id := range conflictID {
			if id == textData.ID {
				check = false
				break
			}
		}
		if check {
			resultData = append(resultData, textData)
		}
		res.ID = textData.ID
		res.Conflict = !check
		response = append(response, res)
	}

	if rows.Err() != nil {
		return result, rows.Err()
	}

	result.resultData = resultData
	result.response = response

	return result, nil
}
