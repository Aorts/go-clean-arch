package bmi

import (
	"context"

	"github.com/bxcodec/go-clean-arch/domain"
)

type BMIRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]domain.BMI, string, error)
	Store(context.Context, *domain.BMI) error
	Delete(ctx context.Context, id int64) error
	GetByName(ctx context.Context, name string) ([]domain.BMI, error)
}

type Service struct {
	bmiRepo BMIRepository
}

func NewService(bmi BMIRepository) *Service {
	return &Service{
		bmiRepo: bmi,
	}
}

func (s *Service) Fetch(ctx context.Context, cursor string, num int64) ([]domain.BMI, string, error) {
	return s.bmiRepo.Fetch(ctx, cursor, num)
}

func (s *Service) GetByName(ctx context.Context, name string) ([]domain.BMI, error) {
	return s.bmiRepo.GetByName(ctx, name)
}

func (s *Service) Store(ctx context.Context, m *domain.BMI) error {
	return s.bmiRepo.Store(ctx, m)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.bmiRepo.Delete(ctx, id)
}
