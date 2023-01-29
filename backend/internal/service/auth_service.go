package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/crypto"
	"github.com/CzarSimon/httputil/jwt"
)

type credentials struct {
	password string
	salt     string
}

func (c credentials) Decode() ([]byte, []byte, error) {
	password, err := fromBase64String(c.password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode password")
	}

	salt, err := fromBase64String(c.salt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode salt")
	}

	return password, salt, nil
}

type AuthService struct {
	UserRepo repository.UserRepository
	Issuer   jwt.Issuer
}

func (s *AuthService) Authenticate(ctx context.Context, req models.LoginRequest) (models.AuthenticatedResponse, error) {
	user, found, err := s.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return models.AuthenticatedResponse{}, err
	}

	if !found {
		return models.AuthenticatedResponse{}, httputil.Unauthorizedf("no user found with the provided email")
	}

	err = verify(ctx, req.Password, getCredentials(user))
	if err != nil {
		return models.AuthenticatedResponse{}, err
	}

	token, err := s.issueToken(ctx, user)
	if err != nil {
		return models.AuthenticatedResponse{}, err
	}

	return models.AuthenticatedResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *AuthService) issueToken(ctx context.Context, user models.User) (string, error) {
	token, err := s.Issuer.Issue(jwt.User{
		ID:    user.ID,
		Roles: []string{models.UserRole},
	}, 24*time.Hour)

	if err != nil {
		return "", fmt.Errorf("failed to issue jwt token: %w", err)
	}

	return token, nil
}

func hash(ctx context.Context, password string) (credentials, error) {
	salt, err := crypto.GenerateAESKey()
	if err != nil {
		return credentials{}, fmt.Errorf("failed to generate salt: %w", err)
	}

	hashtext, err := crypto.DefaultArgon2Hasher().Hash([]byte(password), salt)
	if err != nil {
		return credentials{}, fmt.Errorf("failed to generate salt: %w", err)
	}

	return credentials{
		password: toBase64String(hashtext),
		salt:     toBase64String(salt),
	}, nil
}

func verify(ctx context.Context, password string, creds credentials) error {
	hashtext, salt, err := creds.Decode()
	if err != nil {
		return err
	}

	err = crypto.DefaultArgon2Hasher().Verify([]byte(password), salt, hashtext)
	if err == crypto.ErrHashMissmatch {
		return httputil.UnauthorizedError(err)
	}
	if err != nil {
		return fmt.Errorf("failed to verify password: %w", err)
	}

	return nil
}

func toBase64String(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func fromBase64String(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func getCredentials(user models.User) credentials {
	return credentials{
		password: user.Password,
		salt:     user.Salt,
	}
}
