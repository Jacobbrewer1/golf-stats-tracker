package rounder

import (
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
)

func (r *repository) CreateRound(round *models.Round) error {
	round.Id = 0
	return round.Insert(r.db)
}
