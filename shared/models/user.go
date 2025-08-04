package models

// UserID represents a unique identifier for users.
type UserID string

// UserRegReq contains registration request data.
type UserRegReq struct {
	Username string
	Password string
	Salt     string
}

// UserLoginReq contains authentication request data.
type UserLoginReq struct {
	Username string
	Password string
}

// UserDB represents the persisted user model.
// Contains fields as stored in the database.
type UserDB struct {
	ID       UserID
	Username string
	PassHash string
	Salt     string
}

// User represents the public user model.
// Contains fields returned to clients after authentication.
type User struct {
	ID   UserID
	JWT  string
	Salt string
}
