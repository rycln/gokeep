package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/rycln/gokeep/shared/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// userStorager defines persistence operations for user data
type userStorager interface {
	AddUser(context.Context, *models.UserDB) error
	GetUserByUsername(context.Context, string) (*models.UserDB, error)
}

// passHasher defines password security operations
type passHasher interface {
	Hash(string) (string, error)
	Compare(string, string) error
}

// jwtService defines JWT token operations
type jwtService interface {
	NewJWTString(models.UserID) (string, error)
}

// UserService implements user authentication business logic
type UserService struct {
	strg   userStorager
	hasher passHasher
	jwt    jwtService
}

// NewUserService constructs a new UserService with required dependencies
func NewUserService(strg userStorager, hasher passHasher, jwt jwtService) *UserService {
	return &UserService{
		strg:   strg,
		hasher: hasher,
		jwt:    jwt,
	}
}

// CreateUser handles new user registration:
func (s *UserService) CreateUser(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	hash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	uid := models.UserID(uuid.NewString())

	userDB := &models.UserDB{
		ID:       uid,
		Username: req.Username,
		PassHash: hash,
	}

	err = s.strg.AddUser(ctx, userDB)
	if err != nil {
		return nil, err
	}

	jwt, err := s.jwt.NewJWTString(uid)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:  uid,
		JWT: jwt,
	}, nil
}

// AuthUser handles user authentication:
func (s *UserService) AuthUser(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	userDB, err := s.strg.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = s.hasher.Compare(userDB.PassHash, req.Password)
	if err != nil {
		return nil, err
	}

	jwt, err := s.jwt.NewJWTString(userDB.ID)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:  userDB.ID,
		JWT: jwt,
	}, nil
}
