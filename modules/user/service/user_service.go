// modules/user/service/user_service.go

package service

import (
    "context"
    "errors"

    "github.com/everest/bheri/models"
    "github.com/everest/bheri/modules/user/dto"
    "github.com/everest/bheri/modules/user/repository"
    "github.com/everest/bheri/pkg/helpers"

    "github.com/google/uuid"
)

type UserService interface {
    Register(
        ctx context.Context,
        request dto.RegisterRequest,
    ) (*models.User, string, error) // Updated to return token string

    Login(
        ctx context.Context,
        request dto.LoginRequest,
    ) (*models.User, string, error) // Updated to return token string

    GetByID(
        ctx context.Context,
        id string,
    ) (*models.User, error)
}

type userService struct {
    repository repository.UserRepository
}

func NewUserService(
    repository repository.UserRepository,
) UserService {
    return &userService{
        repository: repository,
    }
}

func (s *userService) Register(
    ctx context.Context,
    request dto.RegisterRequest,
) (*models.User, string, error) {

    if request.Name == "" {
        return nil, "", errors.New("name is required")
    }

    if request.Email == "" {
        return nil, "", errors.New("email is required")
    }

    if len(request.Password) < 8 {
        return nil, "", errors.New("password must be at least 8 characters")
    }

    hash, err := helpers.HashPassword(request.Password)
    if err != nil {
        return nil, "", err
    }

    user, err := s.repository.Create(
        ctx,
        request.Name,
        request.Email,
        hash,
    )
    if err != nil {
        return nil, "", err
    }

    // Generate JWT token for the newly registered user
    token, err := helpers.GenerateToken(user.ID)
    if err != nil {
        return nil, "", err
    }

    return user, token, nil
}

func (s *userService) Login(
    ctx context.Context,
    request dto.LoginRequest,
) (*models.User, string, error) {

    user, passwordHash, err :=
        s.repository.FindByEmail(ctx, request.Email)

    if err != nil {
        return nil, "", err
    }

    if !helpers.CheckPassword(
        request.Password,
        passwordHash,
    ) {
        return nil, "", errors.New("invalid credentials")
    }

    // Generate JWT token upon successful login
    token, err := helpers.GenerateToken(user.ID)
    if err != nil {
        return nil, "", err
    }

    return user, token, nil
}

func (s *userService) GetByID(
    ctx context.Context,
    id string,
) (*models.User, error) {

    if _, err := uuid.Parse(id); err != nil {
        return nil, errors.New("invalid user id")
    }

    return s.repository.FindByID(ctx, id)
}
