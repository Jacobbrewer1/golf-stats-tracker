package rounder

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
)

var (
	ErrNoHolesFound = errors.New("no holes found")
	ErrNoStatsFound = errors.New("no stats found")
)

func (r *repository) CreateHole(hole *models.Hole) error {
	hole.Id = 0
	return hole.Insert(r.db)
}

func (r *repository) GetRoundHoles(roundId int) (*PaginationResponse[models.Hole], error) {
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

	return &PaginationResponse[models.Hole]{
		Items: holes,
		Total: int64(len(holes)),
	}, nil
}

func (r *repository) GetHoleStatsByHoleId(holeId int) (*models.HoleStats, error) {
	sqlStmt := `SELECT id FROM hole_stats WHERE hole_id = ?`

	var id int
	err := r.db.Get(&id, sqlStmt, holeId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoStatsFound
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

func (r *repository) GetAllStatsForPar(userId int, par int64) (*PaginationResponse[HoleWithStats], error) {
	sqlStmt := `
	SELECT
	    r.id AS round_id,
		s.id AS stats_id
	FROM hole h
		INNER JOIN hole_stats s ON h.id = s.hole_id
		INNER JOIN course_details cd ON h.course_details_id = cd.id
		INNER JOIN course c ON cd.course_id = c.id
		INNER JOIN round r ON c.round_id = r.id
	WHERE r.user_id = ?
		AND h.par = ?
	`

	type idStruct struct {
		RoundId int `db:"round_id"`
		StatsId int `db:"stats_id"`
	}

	ids := make([]idStruct, 0)
	err := r.db.Select(&ids, sqlStmt, userId, par)
	if err != nil {
		return nil, fmt.Errorf("failed to get hole stats: %w", err)
	}

	holeStats := make([]*HoleWithStats, 0, len(ids))
	for _, id := range ids {
		round, err := models.RoundById(r.db, id.RoundId)
		if err != nil {
			return nil, fmt.Errorf("failed to get round by ID: %w", err)
		}

		holeStat, err := models.HoleStatsById(r.db, id.StatsId)
		if err != nil {
			return nil, fmt.Errorf("failed to get hole stats by ID: %w", err)
		}

		hole, err := holeStat.GetHole(r.db)
		if err != nil {
			return nil, fmt.Errorf("failed to get hole by ID: %w", err)
		}

		holeStats = append(holeStats, &HoleWithStats{
			Round: round,
			Hole:  hole,
			Stats: holeStat,
		})
	}

	return &PaginationResponse[HoleWithStats]{
		Items: holeStats,
		Total: int64(len(holeStats)),
	}, nil
}

func (r *repository) CountHolesByRoundAndPar(roundId int, par int64) (int, error) {
	sqlStmt := `
	SELECT COUNT(*)
	FROM hole h
		INNER JOIN course_details cd ON h.course_details_id = cd.id
		INNER JOIN course c ON cd.course_id = c.id
	WHERE c.round_id = ?
		AND h.par = ?
	`

	var count int
	err := r.db.Get(&count, sqlStmt, roundId, par)
	if err != nil {
		return 0, fmt.Errorf("failed to count holes: %w", err)
	}

	return count, nil
}

func (r *repository) GetStatsByRoundId(userId int, roundId int) (*PaginationResponse[HoleWithStats], error) {
	sqlStmt := `
	SELECT
		s.id AS stats_id
	FROM hole h
		INNER JOIN hole_stats s ON h.id = s.hole_id
		INNER JOIN course_details cd ON h.course_details_id = cd.id
		INNER JOIN course c ON cd.course_id = c.id
		INNER JOIN round r ON c.round_id = r.id
	WHERE r.user_id = ?
		AND r.id = ?
	`

	ids := make([]int, 0)
	err := r.db.Select(&ids, sqlStmt, userId, roundId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoStatsFound
		default:
			return nil, fmt.Errorf("failed to get hole stats: %w", err)
		}
	}

	round, err := models.RoundById(r.db, roundId)
	if err != nil {
		return nil, fmt.Errorf("failed to get round by ID: %w", err)
	}

	holeStats := make([]*HoleWithStats, 0, len(ids))
	for _, id := range ids {
		holeStat, err := models.HoleStatsById(r.db, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get hole stats by ID: %w", err)
		}

		hole, err := holeStat.GetHole(r.db)
		if err != nil {
			return nil, fmt.Errorf("failed to get hole by ID: %w", err)
		}

		holeStats = append(holeStats, &HoleWithStats{
			Round: round,
			Hole:  hole,
			Stats: holeStat,
		})
	}

	return &PaginationResponse[HoleWithStats]{
		Items: holeStats,
		Total: int64(len(holeStats)),
	}, nil
}

func (r *repository) CountHolesByRoundId(roundId int) (int, error) {
	sqlStmt := `
	SELECT COUNT(*)
	FROM hole h
		INNER JOIN course_details cd ON h.course_details_id = cd.id
		INNER JOIN course c ON cd.course_id = c.id
	WHERE c.round_id = ?
	`

	var count int
	err := r.db.Get(&count, sqlStmt, roundId)
	if err != nil {
		return 0, fmt.Errorf("failed to count holes: %w", err)
	}

	return count, nil
}

func (r *repository) GetStatsByUserId(userId int) (*PaginationResponse[RoundWithStats], error) {
	sqlStmt := `
	SELECT
		rs.id AS stats_id,
		r.id AS round_id,
		c.id AS course_id
	FROM round_stats rs
		INNER JOIN round r ON rs.round_id = r.id
		INNER JOIN course c ON c.round_id = r.id
	WHERE r.user_id = ?
	`

	type idStruct struct {
		StatsId  int `db:"stats_id"`
		RoundId  int `db:"round_id"`
		CourseId int `db:"course_id"`
	}

	ids := make([]idStruct, 0)
	err := r.db.Select(&ids, sqlStmt, userId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoStatsFound
		default:
			return nil, fmt.Errorf("failed to get round stats: %w", err)
		}
	}

	holeStats := make([]*RoundWithStats, 0, len(ids))
	for _, id := range ids {
		round, err := models.RoundById(r.db, id.RoundId)
		if err != nil {
			return nil, fmt.Errorf("failed to get round by ID: %w", err)
		}

		course, err := models.CourseById(r.db, id.CourseId)
		if err != nil {
			return nil, fmt.Errorf("failed to get course by ID: %w", err)
		}

		roundStats, err := models.RoundStatsById(r.db, id.StatsId)
		if err != nil {
			return nil, fmt.Errorf("failed to get round stats by ID: %w", err)
		}

		holeStats = append(holeStats, &RoundWithStats{
			Round:  round,
			Course: course,
			Stats:  roundStats,
		})
	}

	return &PaginationResponse[RoundWithStats]{
		Items: holeStats,
		Total: int64(len(holeStats)),
	}, nil
}
