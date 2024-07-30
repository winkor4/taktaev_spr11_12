package storage

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS users 
	(
		login TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		key TEXT NOT NULL
	);`,
	// `CREATE TABLE IF NOT EXISTS text_data (
	// 	user_login TEXT UNIQUE NOT NULL,
	// 	id TEXT UNIQUE NOT NULL,
	// 	data TEXT NOT NULL
	// );`,
	// `ALTER TABLE text_data ADD UNIQUE (user_login, id)`,
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
)
