package repository

import (
	"database/sql"

	"github.com/snipersune/LetsBetRNM/internal/models"
)

type UserRepository interface {
	CreateUser(username, password string) error
	GetUserByID(id int) (*models.User, error)
	GetUserTeams(id int) ([]models.Team, error)
	GetUserBets(id int) ([]models.Bet, error)
	DeleteUser(id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(username, hashedPassword string) error {
	// Insert user into database
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
	return err
}

func (r *userRepository) DeleteUser(id int) error {
	// Delete user from database
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow("SELECT id, name, createdAt FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserTeams(id int) ([]models.Team, error) {
	rows, err := r.db.Query(
		`SELECT
		teams.id,
		teams.name,
		teams.created_at
		FROM teams
		JOIN team_members ON teams.id = team_members.team_id
		WHERE team_members.user_id = ?;`,
		id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		if err := rows.Scan(&team.ID, &team.Name, &team.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (r *userRepository) GetUserBets(id int) ([]models.Bet, error) {
	queryGameRows, err := r.db.Query(
		`SELECT 
		bets.powerplay_id,
		COALESCE(teams.name, users.username) AS participant,
		bets.placed_at, 
		bets.data
		FROM bets
		LEFT JOIN teams ON bets.team_id = teams.id
		LEFT JOIN users ON bets.user_id = users.id
		WHERE bets.user_id = ? OR bets.team_id IN (
    	SELECT team_id FROM team_members WHERE user_id = ?
		)
		ORDER BY bets.placed_at DESC;
		`, id, id,
	)

	if err != nil {
		return nil, err
	}

	defer queryGameRows.Close()
	return nil, err
}
