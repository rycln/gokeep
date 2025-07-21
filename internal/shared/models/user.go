package models

type UserID string

type UserAuthReq struct {
	Username string
	Password string
}

type UserDB struct {
	ID       UserID
	Username string
	PassHash string
}

type User struct {
	ID  UserID
	JWT string
}
