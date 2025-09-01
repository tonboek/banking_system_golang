package services

import (
	"banking_system_golang/internal/models"
	"banking_system_golang/internal/repository"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, name, username, email, password string) (*models.User, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    return s.repo.Create(ctx, name, username, email, string(hashed))
}

func (s *UserService) Login(ctx context.Context, username, password string) (*models.User, error) {
    user, err := s.repo.FindByUsername(ctx, username)
    if err != nil {
        return nil, errors.New("пользователь не найден")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("неверный пароль")
    }

    token, err := utils.GenerateToken(username)
    if err != nil {
        return nil, errors.New("ошибка при генерации токена")
    }

    return token, nil
}