package service

import (
	"context"
	"fmt"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/timeutil"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.User, error) {
	creds, err := hash(ctx, req.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user := req.User(creds.password, creds.salt)
	err = s.UserRepo.Save(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	user, exists, err := s.UserRepo.Find(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	if !exists {
		return models.User{}, httputil.NotFoundError(err)
	}

	return user, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.UserRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) ChangePassword(ctx context.Context, req models.ChangePasswordRequest) error {
	user, exits, err := s.UserRepo.Find(ctx, req.UserID)
	if err != nil {
		return err
	}
	if !exits {
		return httputil.PreconditionRequiredError(err)
	}

	err = verify(ctx, req.OldPassword, getCredentials(user))
	if err != nil {
		return err
	}

	creds, err := hash(ctx, req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = creds.password
	user.Salt = creds.salt
	user.UpdatedAt = timeutil.Now()

	return s.UserRepo.Save(ctx, user)
}
