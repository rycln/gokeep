package storage

const sqlAddUser = `
	INSERT INTO users (id, username, password_hash, salt) 
	VALUES ($1, $2, $3, $4) 
`

const sqlGetUserByUsername = `
	SELECT 
		id, 
		username, 
		password_hash,
		salt 
	FROM users 
	WHERE username = $1
`
