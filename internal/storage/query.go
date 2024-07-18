package storage

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS users 
	(
		login TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS text_data (
		id TEXT NOT NULL,
		user_login TEXT NOT NULL,
		data TEXT NOT NULL
	);`,
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
)
