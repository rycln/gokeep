package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/rycln/gokeep/internal/shared/models"
)

type userStorager interface {
	AddUser(context.Context, *models.UserDB) error
	GetUserByUsername(context.Context, string) (*models.UserDB, error)
}

type passHasher interface {
	Hash(string) (string, error)
	Compare(string, string) error
}

type jwtService interface {
	NewJWTString(models.UserID) (string, error)
}

type UserService struct {
	strg   userStorager
	hasher passHasher
	jwt    jwtService
}

func NewUserService(strg userStorager, hasher passHasher, jwt jwtService) *UserService {
	return &UserService{
		strg:   strg,
		hasher: hasher,
		jwt:    jwt,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	hash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	uid := models.UserID(uuid.NewString())

	userDB := &models.UserDB{
		ID:       models.UserID(uuid.NewString()),
		Username: req.Username,
		PassHash: hash,
	}
	err = s.strg.AddUser(ctx, userDB)
	if err != nil {
		return nil, err
	}

	//добавить refresh jwt
	jwt, err := s.jwt.NewJWTString(uid)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:  uid,
		JWT: jwt,
	}

	return user, nil
}

func (s *UserService) AuthUser(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	userDB, err := s.strg.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	err = s.hasher.Compare(userDB.PassHash, req.Password)
	if err != nil {
		return nil, err
	}

	//добавить refresh jwt
	jwt, err := s.jwt.NewJWTString(userDB.ID)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:  userDB.ID,
		JWT: jwt,
	}

	return user, nil
}
