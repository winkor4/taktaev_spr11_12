package storage

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS users 
	(
		login TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS text_data (
		user_login TEXT UNIQUE NOT NULL,
		id TEXT UNIQUE NOT NULL,
		data TEXT NOT NULL
	);`,
	`ALTER TABLE text_data ADD UNIQUE (user_login, id)`,
}

var (
	queryRegister = `
	INSERT INTO users 
	(
		login, 
		password
	)
	VALUES 
	(
		$1, 
		$2
	)
	ON CONFLICT (login) DO NOTHING;`

	queryPassword = `
	SELECT
		password
	FROM
		users
	WHERE 
		login = $1`

	queryUploadTextData = `
	INSERT INTO text_data
	(
		user_login,
		id,
		data
	)
	VALUES 
	(
		$1,
		$2,
		$3
	)
	ON CONFLICT (user_login, id) DO NOTHING;`

	queryConflictID = `
	SELECT
		id
	FROM
		text_data
	WHERE
		user_login = $1
		AND id IN ($2)`
)
