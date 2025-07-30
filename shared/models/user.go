package models

// UserID represents a unique identifier for users.
type UserID string

// UserAuthReq contains authentication request data.
// Used for both login and registration operations.
type UserAuthReq struct {
	Username string
	Password string
}

// UserDB represents the persisted user model.
// Contains fields as stored in the database.
type UserDB struct {
	ID       UserID
	Username string
	PassHash string
}

// User represents the public user model.
// Contains fields returned to clients after authentication.
type User struct {
	ID  UserID
	JWT string
}
