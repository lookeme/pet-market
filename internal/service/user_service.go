package service

import (
	"context"
	"pet-market/api"
	"pet-market/internal/logger"
	"pet-market/internal/repository"
	"pet-market/internal/security"

	"go.uber.org/zap"
)

type UsrService struct {
	userRepository repository.UserRepository
	auth           security.Authorization
	Log            *logger.Logger
}

func NewUserService(userRepository repository.UserRepository, auth security.Authorization, log *logger.Logger) *UsrService {
	return &UsrService{
		userRepository,
		auth,
		log,
	}
}

func (s *UsrService) CreateUser(ctx context.Context, usr api.User) error {
	hash, err := s.auth.HashPassword(usr.Password)
	if err != nil {
		return err
	}
	_, err = s.userRepository.Save(ctx, usr.Login, hash)
	s.Log.Log.Info("user created", zap.String("login ", usr.Login))
	if err != nil {
		return err
	}
	return nil
}

func (s *UsrService) GetUserByName(ctx context.Context, userName string) (*api.User, error) {
	usr, err := s.userRepository.GetUserByLogin(ctx, userName)
	return &usr, err
}
