package repository

import (
	"database/sql"

	"github.com/snipersune/LetsBetRNM/internal/models"
)

type TeamRepository interface {
	CreateTeam(name string) error
	AddUserToTeam(user_id, team_id int) error
	RemoveUserFromTeam(user_id, team_id int) error
	GetTeamByID(id int) (*models.Team, error)
	GetTeamMembers(id int) ([]models.User, error)
	DeleteTeam(id int) error
}

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) CreateTeam(name string) error {
	// Insert team into database
	_, err := r.db.Exec("INSERT INTO teams VALUES ?", name)
	return err
}

func (r *teamRepository) DeleteTeam(id int) error {
	// Delete team from database
	_, err := r.db.Exec("DELETE FROM teams WHERE id = $1", id)
	return err
}

func (r *teamRepository) GetTeamByID(id int) (*models.Team, error) {
	team := &models.Team{}
	err := r.db.QueryRow("SELECT id, name, createdAt FROM teams WHERE id = $1", id).
		Scan(&team.ID, &team.Name, &team.CreatedAt)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (r *teamRepository) AddUserToTeam(user_id, team_id int) error {
	_, err := r.db.Exec("INSERT INTO teams_members (user_id, team_id) VALUES (?, ?)", user_id, team_id)
	return err
}

func (r *teamRepository) RemoveUserFromTeam(user_id, team_id int) error {
	_, err := r.db.Exec("DELETE FROM teams_members where (user_id, team_id) VALUES (?, ?)", user_id, team_id)
	return err
}

func (r *teamRepository) GetTeamMembers(id int) ([]models.User, error) {
	rows, err := r.db.Query(
		`SELECT
		users.id,
		users.name,
		users.created_at
		FROM users
		JOIN team_members ON users.id = team_members.user_id
		WHERE team_members.user_id = ?;`,
		id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var team models.User
		if err := rows.Scan(&team.ID, &team.Name, &team.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, team)
	}

	return users, nil
}
