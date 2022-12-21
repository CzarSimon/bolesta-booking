package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/crypto"
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
