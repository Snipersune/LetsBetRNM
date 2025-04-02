package service

import (
	"errors"

	"github.com/snipersune/LetsBetRNM/internal/models"
	"github.com/snipersune/LetsBetRNM/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(username, password string) error
	DeleteUser(id int) error
	GetUserByID(id int) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// RegisterUser validates input and creates a new user
func (s *userService) RegisterUser(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password required")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(username, string(hashedPassword))
}

func (s *userService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetUserByID(id)
}
