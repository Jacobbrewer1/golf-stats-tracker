package rounder

import (
	"fmt"

	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
)

func (r *repository) SaveRoundHitStats(roundHitStats ...*models.RoundHitStats) error {
	for _, rhs := range roundHitStats {
		err := rhs.SaveOrUpdate(r.db)
		if err != nil {
			return fmt.Errorf("failed to save round hit stats: %w", err)
		}
	}

	return nil
}

func (r *repository) GetRoundHitStatsByRoundStatsId(roundStatsId int) (*PaginationResponse[models.RoundHitStats], error) {
	sqlStmt := `SELECT id FROM round_hit_stats WHERE round_stats_id = ?`

	ids := make([]int, 0)
	err := r.db.Select(&ids, sqlStmt, roundStatsId)
	if err != nil {
		return nil, fmt.Errorf("failed to get round hit stats: %w", err)
	}

	roundHitStats := make([]*models.RoundHitStats, 0, len(ids))
	for _, id := range ids {
		rhs, err := models.RoundHitStatsById(r.db, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get round hit stats by ID: %w", err)
		}
		roundHitStats = append(roundHitStats, rhs)
	}

	return &PaginationResponse[models.RoundHitStats]{
		Items: roundHitStats,
		Total: int64(len(roundHitStats)),
	}, nil
}

func (r *repository) GetRoundStatsByRoundId(roundId int) (*models.RoundStats, error) {
	sqlStmt := `SELECT id FROM round_stats WHERE round_id = ?`

	var id int
	err := r.db.Get(&id, sqlStmt, roundId)
	if err != nil {
		return nil, fmt.Errorf("failed to get round stats: %w", err)
	}

	return models.RoundStatsById(r.db, id)
}

func (r *repository) GetRoundHitStats(roundId int) (*PaginationResponse[models.RoundHitStats], error) {
	sqlStmt := `SELECT id FROM round_hit_stats WHERE round_id = ?`

	ids := make([]int, 0)
	err := r.db.Select(&ids, sqlStmt, roundId)
	if err != nil {
		return nil, fmt.Errorf("failed to get round hit stats: %w", err)
	}

	roundHitStats := make([]*models.RoundHitStats, 0, len(ids))
	for _, id := range ids {
		rhs, err := models.RoundHitStatsById(r.db, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get round hit stats by ID: %w", err)
		}
		roundHitStats = append(roundHitStats, rhs)
	}

	return &PaginationResponse[models.RoundHitStats]{
		Items: roundHitStats,
		Total: int64(len(roundHitStats)),
	}, nil
}

func (r *repository) GetUserHitStats(userId int) (*PaginationResponse[models.RoundHitStats], error) {
	sqlStmt := `
		SELECT rhs.id
		FROM round_hit_stats rhs
			JOIN round_stats rs ON rhs.round_stats_id = rs.id
			JOIN round r ON rs.round_id = r.id
		WHERE r.user_id = ?`

	ids := make([]int, 0)
	err := r.db.Select(&ids, sqlStmt, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user hit stats: %w", err)
	}

	roundHitStats := make([]*models.RoundHitStats, 0, len(ids))
	for _, id := range ids {
		rhs, err := models.RoundHitStatsById(r.db, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get round hit stats by ID: %w", err)
		}
		roundHitStats = append(roundHitStats, rhs)
	}

	return &PaginationResponse[models.RoundHitStats]{
		Items: roundHitStats,
		Total: int64(len(roundHitStats)),
	}, nil
}
