package service

import (
	"errors"

	"github.com/snipersune/LetsBetRNM/internal/models"
	"github.com/snipersune/LetsBetRNM/internal/repository"
)

type TeamService interface {
	RegisterTeam(name string) error
	DeleteTeam(id int) error
	GetTeamByID(id int) (*models.Team, error)
}

type teamService struct {
	repo repository.TeamRepository
}

// NewUserService creates a new instance of UserService
func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamService{repo: repo}
}

// RegisterUser validates input and creates a new user
func (s *teamService) RegisterTeam(name string) error {
	if name == "" {
		return errors.New("team name required")
	}

	return s.repo.CreateTeam(name)
}

func (s *teamService) DeleteTeam(id int) error {
	return s.repo.DeleteTeam(id)
}

// GetTeamByID retrieves a Team by ID
func (s *teamService) GetTeamByID(id int) (*models.Team, error) {
	return s.repo.GetTeamByID(id)
}
