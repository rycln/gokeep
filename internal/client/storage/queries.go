package storage

const sqlCreateItemsTable = `
	CREATE TABLE IF NOT EXISTS items (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		encrypt_content BLOB NOT NULL,
		metadata TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
`
const sqlAddItem = `
	INSERT INTO items
	(id, user_id, type, name, encrypt_content, metadata) 
	VALUES ($1, $2, $3, $4, $5, $6)
`

const sqlGetItemByID = `
	SELECT 
		encrypt_content
	FROM items
	WHERE id = $1
`

const sqlGetUserItemsInfo = `
	SELECT 
		id,
		user_id,
		type,
		name, 
		metadata,
		updated_at 
	FROM items
	WHERE user_id = $1
`
