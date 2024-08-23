package rounder

import "github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"

func (r *repository) CreateHole(hole *models.Hole) error {
	hole.Id = 0
	return hole.Insert(r.db)
}
