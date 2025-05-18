package app

import (
	"context"
	"github.com/invinciblewest/gophermart/internal/helper"
	"github.com/invinciblewest/gophermart/internal/model"
	"github.com/invinciblewest/gophermart/internal/repository"
	"github.com/invinciblewest/gophermart/internal/usecase"
)

type UserUseCase struct {
	userRepository repository.UserRepository
	authUseCase    usecase.AuthUseCase
}

func NewUserUseCase(userRepository repository.UserRepository, authUseCase usecase.AuthUseCase) *UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
		authUseCase:    authUseCase,
	}
}

func (us *UserUseCase) RegisterAndLogin(ctx context.Context, user *model.User) (string, error) {
	if user.Login == "" || user.Password == "" {
		return "", helper.ErrEmptyLoginOrPassword
	}

	user.Password = us.authUseCase.HashPassword(user.Password)

	if err := us.userRepository.CreateUser(ctx, user); err != nil {
		return "", err
	}

	return us.authUseCase.GenerateToken(user.ID)
}

func (us *UserUseCase) Login(ctx context.Context, user model.User) (string, error) {
	if user.Login == "" || user.Password == "" {
		return "", helper.ErrEmptyLoginOrPassword
	}

	receivedUser, err := us.userRepository.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return "", err
	}

	if !us.authUseCase.VerifyPassword(receivedUser, user.Password) {
		return "", helper.ErrInvalidPassword
	}

	return us.authUseCase.GenerateToken(receivedUser.ID)
}
