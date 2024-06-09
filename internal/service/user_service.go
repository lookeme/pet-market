package service

import (
	"pet-market/api"
	"pet-market/internal/logger"
	"pet-market/internal/repository"
	"pet-market/internal/security"

	"go.uber.org/zap"
)

type UsrService struct {
	userRepository repository.UserRepository
	auth           *security.Authorization
	Log            *logger.Logger
}

func NewUserService(userRepository repository.UserRepository, auth *security.Authorization, log *logger.Logger) *UsrService {
	return &UsrService{
		userRepository,
		auth,
		log,
	}
}

func (s *UsrService) CreateUser(usr api.User) error {
	hash, err := s.auth.HashPassword(usr.Password)
	if err != nil {
		return err
	}
	_, err = s.userRepository.Save(usr.Login, hash)
	s.Log.Log.Info("user created", zap.String("login ", usr.Login))
	if err != nil {
		return err
	}
	return nil
}

func (s *UsrService) GetUserByName(userName string) (*api.User, error) {
	usr, err := s.userRepository.GetUserByName(userName)
	return &usr, err
}
