// modules/api_service/service/api_service_service.go

package service

import (
	"context"
	"errors"

	"github.com/everest/bheri/models"
	"github.com/everest/bheri/modules/api_service/dto"
	"github.com/everest/bheri/modules/api_service/repository"
)

type APIServiceService interface {
	Create(
		ctx context.Context,
		request dto.CreateAPIServiceRequest,
	) (*models.APIService, error)

	Get(
		ctx context.Context,
		id string,
	) (*models.APIService, error)

	List(
		ctx context.Context,
	) ([]models.APIService, error)

	ListActive(
		ctx context.Context,
	) ([]models.APIService, error)

	Update(
		ctx context.Context,
		id string,
		request dto.UpdateAPIServiceRequest,
	) (*models.APIService, error)

	Delete(
		ctx context.Context,
		id string,
	) error
}

type apiServiceService struct {
	repository repository.APIServiceRepository
}

func NewAPIServiceService(
	repository repository.APIServiceRepository,
) APIServiceService {
	return &apiServiceService{
		repository: repository,
	}
}

func (s *apiServiceService) Create(
	ctx context.Context,
	request dto.CreateAPIServiceRequest,
) (*models.APIService, error) {

	if request.Name == "" {
		return nil, errors.New(
			"service name is required",
		)
	}

	if request.Endpoint == "" {
		return nil, errors.New(
			"endpoint is required",
		)
	}

	if request.PricePerRequest <= 0 {
		return nil, errors.New(
			"price must be greater than zero",
		)
	}

	if request.Asset == "" {
		request.Asset = "USDC"
	}

	if request.Network == "" {
		request.Network = "solana-devnet"
	}

	return s.repository.Create(
		ctx,
		request,
	)
}

func (s *apiServiceService) Get(
	ctx context.Context,
	id string,
) (*models.APIService, error) {

	return s.repository.FindByID(
		ctx,
		id,
	)
}

func (s *apiServiceService) List(
	ctx context.Context,
) ([]models.APIService, error) {

	return s.repository.FindAll(ctx)
}

func (s *apiServiceService) ListActive(
	ctx context.Context,
) ([]models.APIService, error) {

	return s.repository.FindActive(ctx)
}

func (s *apiServiceService) Update(
	ctx context.Context,
	id string,
	request dto.UpdateAPIServiceRequest,
) (*models.APIService, error) {

	if request.Name == "" {
		return nil, errors.New(
			"service name is required",
		)
	}

	if request.Endpoint == "" {
		return nil, errors.New(
			"endpoint is required",
		)
	}

	if request.PricePerRequest <= 0 {
		return nil, errors.New(
			"price must be greater than zero",
		)
	}

	return s.repository.Update(
		ctx,
		id,
		request,
	)
}

func (s *apiServiceService) Delete(
	ctx context.Context,
	id string,
) error {

	return s.repository.Delete(
		ctx,
		id,
	)
}
