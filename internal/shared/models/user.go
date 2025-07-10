package models

type UserID string

type UserAuth struct {
	Login    string
	Password string
}

type UserDB struct {
	ID       UserID
	Login    string
	PassHash string
}
