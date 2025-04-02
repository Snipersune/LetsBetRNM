package repository

import (
	"database/sql"

	"github.com/snipersune/LetsBetRNM/internal/models"
)

type PowerplayRepository interface {
	CreatePowerplay(pplay *models.Powerplay) error
	DeletePowerplay(id int) error
	UpdatePowerplay(id int, pplay *models.Powerplay) error
	GetPowerplayByID(id int) (*models.Powerplay, error)
}

type powerplayRepository struct {
	db *sql.DB
}

func NewPowerplayRepository(db *sql.DB) PowerplayRepository {
	return &powerplayRepository{db: db}
}

func (r *powerplayRepository) CreatePowerplay(pplay *models.Powerplay) error {
	_, err := r.db.Exec(`
		INSERT INTO powerplays 
		(
		home_teams, 
		away_teams,
		percs,
		odds,
		) 
		VALUES (?, ?, ?, ?)`,
		pplay.HomeTeams,
		pplay.AwayTeams,
		pplay.Percentages,
		pplay.Odds,
	)
	return err
}

func (r *powerplayRepository) DeletePowerplay(id int) error {
	_, err := r.db.Exec("DELETE FROM powerplays WHERE id = $1", id)
	return err
}

func (r *powerplayRepository) UpdatePowerplay(id int, pplay *models.Powerplay) error {
	_, err := r.db.Exec("UPDATE powerplays SET percs = $1, odds = $2 WHERE id = $3", pplay.Percentages, pplay.Odds, pplay.ID)
	return err
}

func (r *powerplayRepository) GetPowerplayByID(id int) (*models.Powerplay, error) {
	pplay := &models.Powerplay{}
	err := r.db.QueryRow("SELECT id, home_teams, away_teams, percs, odds FROM powerplays WHERE id = $1", id).
		Scan(&pplay.ID, &pplay.HomeTeams, &pplay.AwayTeams, &pplay.Percentages, &pplay.Odds)
	if err != nil {
		return nil, err
	}
	return pplay, nil
}
