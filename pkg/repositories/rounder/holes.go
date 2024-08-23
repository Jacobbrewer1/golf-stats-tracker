package rounder

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
)

var (
	ErrNoHolesFound     = errors.New("no holes found")
	ErrNoHoleStatsFound = errors.New("no hole stats found")
)

func (r *repository) CreateHole(hole *models.Hole) error {
	hole.Id = 0
	return hole.Insert(r.db)
}

func (r *repository) GetRoundHoles(roundId int) ([]*models.Hole, error) {
	sqlStmt := `
	SELECT h.id
	FROM hole h
		INNER JOIN course_details cd ON h.course_details_id = cd.id
		INNER JOIN course c ON cd.course_id = c.id
	WHERE c.round_id = ?
	`

	holeIds := make([]int, 0)
	err := r.db.Select(&holeIds, sqlStmt, roundId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoHolesFound
		default:
			return nil, fmt.Errorf("failed to get hole IDs: %w", err)
		}
	}

	holes := make([]*models.Hole, 0, len(holeIds))
	for _, id := range holeIds {
		h, err := models.HoleById(r.db, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get hole by ID: %w", err)
		}
		holes = append(holes, h)
	}

	return holes, nil
}

func (r *repository) GetHoleStatsByHoleId(holeId int) (*models.HoleStats, error) {
	sqlStmt := `SELECT id FROM hole_stats WHERE hole_id = ?`

	var id int
	err := r.db.Get(&id, sqlStmt, holeId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoHoleStatsFound
		default:
			return nil, fmt.Errorf("failed to get hole stats by hole ID: %w", err)
		}
	}

	return models.HoleStatsById(r.db, id)
}

func (r *repository) GetHoleById(id int) (*models.Hole, error) {
	return models.HoleById(r.db, id)
}

func (r *repository) SaveHoleStats(holeStats *models.HoleStats) error {
	err := holeStats.SaveOrUpdate(r.db)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoAffectedRows):
			break
		default:
			return fmt.Errorf("failed to save hole stats: %w", err)
		}
	}

	return nil
}
