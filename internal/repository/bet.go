package repository

import "github.com/snipersune/LetsBetRNM/cmd/models"

type BetRepository interface {
	CreateBet(*models.Bet) error
	GetBetByID(id int) (*models.Bet, error)
	DeleteBet(id int) error
}
