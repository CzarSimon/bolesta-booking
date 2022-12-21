package service

import (
	"context"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/httputil"
)

type CabinService struct {
	CabinRepo repository.CabinRepository
}

func (s *CabinService) GetCabin(ctx context.Context, id string) (models.Cabin, error) {
	cabin, exists, err := s.CabinRepo.Find(ctx, id)
	if err != nil {
		return models.Cabin{}, err
	}

	if !exists {
		return models.Cabin{}, httputil.NotFoundError(err)
	}

	return cabin, nil
}

func (s *CabinService) GetCabins(ctx context.Context) ([]models.Cabin, error) {
	cabins, err := s.CabinRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return cabins, nil
}
