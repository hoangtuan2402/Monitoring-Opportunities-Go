package service

import (
	"Monitoring-Opportunities/src/dto"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrCreateUserValidate = errors.New("failed to create a user")
	ErrUpdateUserVaildate = errors.New("failed to update a user")
	ErrDeleteUuser        = errors.New("failed to delete a user")
	ErrDuplicateUserEmail = errors.New("user with that email is already exists")
	ErrDuplicateUsername  = errors.New("user with that username is already exists")
)

type UserService interface {
	GetAll() ([]dto.UserDTO, error)
	Create(user dto.CreateUser) (dto.UserDTO, error)
	Update(user dto.UpdateUser, userID uuid.UUID) (dto.UserDTO, error)
	Delete(id uuid.UUID) (dto.UserDTO, error)
	FindByID(id uuid.UUID) (dto.UserDTO, error)
	FindByEmail(email string) (dto.UserDTO, error)
}

type userService struct {
}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) GetAll() ([]dto.UserDTO, error) {

	return []dto.UserDTO{}, nil
}

func (s *userService) Create(form dto.CreateUser) (dto.UserDTO, error) {

	return dto.UserDTO{
		UUID:     uuid.New(),
		Username: form.Username,
		Email:    form.Email,
	}, nil
}

func (s *userService) Update(user dto.UpdateUser, userID uuid.UUID) (dto.UserDTO, error) {

	return dto.UserDTO{
		UUID:     uuid.New(),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *userService) Delete(id uuid.UUID) (dto.UserDTO, error) {

	return dto.UserDTO{
		UUID:     id,
		Username: "username deleted by id",
		Email:    "email username deleted by id",
	}, nil
}

func (s *userService) FindByID(id uuid.UUID) (dto.UserDTO, error) {

	return dto.UserDTO{
		UUID:     id,
		Username: "username find by id",
		Email:    "email username find by id",
	}, nil
}

func (s *userService) FindByEmail(email string) (dto.UserDTO, error) {

	return dto.UserDTO{
		UUID:     uuid.New(),
		Username: "User Email",
		Email:    email,
	}, nil
}
