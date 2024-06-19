package service

import (
	"context"
	"pet-market/api"
	"pet-market/internal/logger"
	"pet-market/internal/repository"
	"pet-market/internal/security"

	"go.uber.org/zap"
)

type UserService struct {
	userRepository repository.IUserRepository
	auth           security.Authorization
	Log            *logger.Logger
}

func NewUserService(userRepository repository.IUserRepository, auth security.Authorization, log *logger.Logger) *UserService {
	return &UserService{
		userRepository,
		auth,
		log,
	}
}

func (s *UserService) CreateUser(ctx context.Context, usr api.User) (int, error) {
	hash, err := security.HashPassword(usr.Password)
	if err != nil {
		return 0, err
	}
	ID, err := s.userRepository.Save(ctx, usr.Login, hash)
	s.Log.Log.Info("user created", zap.String("login ", usr.Login))
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (s *UserService) GetUserByName(ctx context.Context, userName string) (*api.User, error) {
	usr, err := s.userRepository.GetUserByLogin(ctx, userName)
	return &usr, err
}
