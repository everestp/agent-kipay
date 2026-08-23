// modules/api_service/repository/api_service_repository.go

package repository

import (
	"context"

	"github.com/everest/bheri/models"
	"github.com/everest/bheri/modules/api_service/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServiceRepository interface {
	Create(
		ctx context.Context,
		request dto.CreateAPIServiceRequest,
	) (*models.APIService, error)

	FindByID(
		ctx context.Context,
		id string,
	) (*models.APIService, error)

	FindAll(
		ctx context.Context,
	) ([]models.APIService, error)

	FindActive(
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

type apiServiceRepository struct {
	db *pgxpool.Pool
}

func NewAPIServiceRepository(
	db *pgxpool.Pool,
) APIServiceRepository {
	return &apiServiceRepository{
		db: db,
	}
}

func (r *apiServiceRepository) Create(
	ctx context.Context,
	request dto.CreateAPIServiceRequest,
) (*models.APIService, error) {

	var service models.APIService

err := r.db.QueryRow(
    ctx,
    `
    INSERT INTO api_services (
        name,
        category,
        endpoint,
        price_per_request,
        asset,
        network,
        description
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING
        id,
        name,
        COALESCE(category, ''),
        endpoint,
        price_per_request,
        asset,
        network,
        COALESCE(description, ''),
        provider_reputation,
        status,
        created_at,
        updated_at
    `,
    request.Name,
    request.Category,
    request.Endpoint,
    request.PricePerRequest,
    request.Asset,
    request.Network,
    request.Description,
).Scan(
    &service.ID,
    &service.Name,
    &service.Category,
    &service.Endpoint,
    &service.PricePerRequest,
    &service.Asset,
    &service.Network,
    &service.Description,
    &service.ProviderReputation,
    &service.Status,
    &service.CreatedAt,
    &service.UpdatedAt,
)

	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *apiServiceRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.APIService, error) {

	var service models.APIService

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			name,
			COALESCE(category, ''),
			endpoint,
			price_per_request,
			asset,
			network,
			COALESCE(description, ''),
			provider_reputation,
			active,
			created_at,
			updated_at
		FROM api_services
		WHERE id = $1
		`,
		id,
	).Scan(
		&service.ID,
		&service.Name,
		&service.Category,
		&service.Endpoint,
		&service.PricePerRequest,
		&service.Asset,
		&service.Network,
		&service.Description,
		&service.ProviderReputation,
		&service.Active,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *apiServiceRepository) FindAll(
	ctx context.Context,
) ([]models.APIService, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			name,
			COALESCE(category, ''),
			endpoint,
			price_per_request,
			asset,
			network,
			COALESCE(description, ''),
			provider_reputation,
			active,
			created_at,
			updated_at
		FROM api_services
		ORDER BY created_at DESC
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	services := make([]models.APIService, 0)

	for rows.Next() {

		var service models.APIService

		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Category,
			&service.Endpoint,
			&service.PricePerRequest,
			&service.Asset,
			&service.Network,
			&service.Description,
			&service.ProviderReputation,
			&service.Active,
			&service.CreatedAt,
			&service.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		services = append(
			services,
			service,
		)
	}

	return services, rows.Err()
}

func (r *apiServiceRepository) FindActive(
	ctx context.Context,
) ([]models.APIService, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			name,
			COALESCE(category, ''),
			endpoint,
			price_per_request,
			asset,
			network,
			COALESCE(description, ''),
			provider_reputation,
			active,
			created_at,
			updated_at
		FROM api_services
		WHERE active = TRUE
		ORDER BY provider_reputation DESC
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	services := make([]models.APIService, 0)

	for rows.Next() {

		var service models.APIService

		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Category,
			&service.Endpoint,
			&service.PricePerRequest,
			&service.Asset,
			&service.Network,
			&service.Description,
			&service.ProviderReputation,
			&service.Active,
			&service.CreatedAt,
			&service.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		services = append(
			services,
			service,
		)
	}

	return services, rows.Err()
}

func (r *apiServiceRepository) Update(
	ctx context.Context,
	id string,
	request dto.UpdateAPIServiceRequest,
) (*models.APIService, error) {

	var service models.APIService

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE api_services
		SET
			name = $1,
			category = $2,
			endpoint = $3,
			price_per_request = $4,
			description = $5,
			active = $6,
			updated_at = NOW()
		WHERE id = $7
		RETURNING
			id,
			name,
			COALESCE(category, ''),
			endpoint,
			price_per_request,
			asset,
			network,
			COALESCE(description, ''),
			provider_reputation,
			active,
			created_at,
			updated_at
		`,
		request.Name,
		request.Category,
		request.Endpoint,
		request.PricePerRequest,
		request.Description,
		request.Active,
		id,
	).Scan(
		&service.ID,
		&service.Name,
		&service.Category,
		&service.Endpoint,
		&service.PricePerRequest,
		&service.Asset,
		&service.Network,
		&service.Description,
		&service.ProviderReputation,
		&service.Active,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *apiServiceRepository) Delete(
	ctx context.Context,
	id string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		DELETE FROM api_services
		WHERE id = $1
		`,
		id,
	)

	return err
}
