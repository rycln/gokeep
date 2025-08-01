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

const sqlDeleteItem = `
	UPDATE items 
	SET deleted = true, 
		updated_at = $1 
	WHERE id = $2 AND user_id = $3
`

const sqlAddItem = `
	INSERT INTO items (id, user_id, type, name, metadata, data, updated_at, is_deleted)
    VALUES ($1, $2, $3, $4, $5, $6, $7, false)
	ON CONFLICT (id) DO 
		UPDATE 
        SET type = $3, 
			name = $4, 
			metadata = $5, 
			data = $6, 
			updated_at = $7, 
			is_deleted = false`

const sqlGetUserItems = `
	SELECT 
		id, 
		type, 
		name, 
		metadata, 
		data, 
		updated_at, 
		is_deleted 
	FROM items 
	WHERE user_id = $1
`
