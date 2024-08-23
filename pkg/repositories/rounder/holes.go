package rounder

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
)

var (
	ErrNoHolesFound = errors.New("no holes found")
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
