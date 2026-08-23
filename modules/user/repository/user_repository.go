// modules/user/repository/user_repository.go

package repository

import (
	"context"

	"github.com/everest/bheri/models"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		name string,
		email string,
		passwordHash string,
	) (*models.User, error)

	FindByEmail(
		ctx context.Context,
		email string,
	) (*models.User, string, error)

	FindByID(
		ctx context.Context,
		id string,
	) (*models.User, error)
}

type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
) (*models.User, error) {

	user := &models.User{}

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO users (
			name,
			email,
			password_hash
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			name,
			email,
			status,
			created_at,
			updated_at
		`,
		name,
		email,
		passwordHash,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, string, error) {

	user := &models.User{}
	var passwordHash string

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			name,
			email,
			password_hash,
			status,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
		`,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&passwordHash,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, "", err
	}

	return user, passwordHash, nil
}

func (r *userRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.User, error) {

	user := &models.User{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			name,
			email,
			status,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
