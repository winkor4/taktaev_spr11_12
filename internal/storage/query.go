package storage

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS users 
	(
		login TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		key TEXT NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS content (
		id TEXT UNIQUE NOT NULL,
		user_login TEXT NOT NULL,
		name TEXT NOT NULL,
		data TEXT NOT NULL,
		data_key TEXT NOT NULL
	);`,
}

var (
	queryAddUser = `
	INSERT INTO users 
	(
		login, 
		password,
		key
	)
	VALUES 
	(
		$1, 
		$2,
		$3
	)
	ON CONFLICT (login) DO NOTHING;`

	queryPassword = `
	SELECT
		password
	FROM
		users
	WHERE 
		login = $1`

	queryKey = `
	SELECT
		key
	FROM
		users
	WHERE 
		login = $1`

	queryInsertContent = `
	INSERT INTO content
	(
		id,
		user_login,
		name,
		data,
		data_key
	)
	VALUES 
	(
		$1,
		$2,
		$3,
		$4,
		$5
	)`

	queryGetContent = `
	SELECT 
		content.data,
		content.data_key,
		users.key
	FROM 
		content as content
		LEFT JOIN users as users
		ON content.user_login = users.login
	WHERE
		user_login = $1
		AND name = $2`
)
