package storage

const sqlAddUser = `
	INSERT INTO users (id, username, password_hash) 
	VALUES ($1, $2, $3) 
`

const sqlGetUserByUsername = `
	SELECT 
		id, 
		username, 
		password_hash 
	FROM users 
	WHERE username = $1
`
