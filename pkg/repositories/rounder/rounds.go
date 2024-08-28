package rounder

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
)

func (r *repository) CreateRound(round *models.Round) error {
	round.Id = 0
	return round.Insert(r.db)
}

func (r *repository) GetRoundById(id int) (*models.Round, error) {
	return models.RoundById(r.db, id)
}

func (r *repository) GetRoundDetailsByRoundId(roundId int) (*RoundDetails, error) {
	round, err := r.GetRoundById(roundId)
	if err != nil {
		return nil, err
	}

	sqlStmt := `SELECT id FROM course WHERE round_id = ?`

	var courseID int
	err = r.db.Get(&courseID, sqlStmt, roundId)
	if err != nil {
		return nil, fmt.Errorf("failed to get course ID: %w", err)
	}

	course, err := models.CourseById(r.db, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by ID: %w", err)
	}

	sqlStmt = `SELECT id FROM course_details WHERE course_id = ?`

	var courseDetailsID int
	err = r.db.Get(&courseDetailsID, sqlStmt, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course details ID: %w", err)
	}

	courseDetails, err := models.CourseDetailsById(r.db, courseDetailsID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course details by ID: %w", err)
	}

	sqlStmt = `SELECT id FROM hole WHERE course_details_id = ?`

	var holeIDs []int
	err = r.db.Select(&holeIDs, sqlStmt, courseDetailsID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hole IDs: %w", err)
	}

	holes := make([]*models.Hole, 0, len(holeIDs))
	for _, id := range holeIDs {
		h, err := models.HoleById(r.db, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get hole by ID: %w", err)
		}
		holes = append(holes, h)
	}

	return &RoundDetails{
		Round:         round,
		Course:        course,
		CourseDetails: courseDetails,
		Holes:         holes,
	}, nil
}

func (r *repository) GetRoundsByUserId(userId int) ([]*models.Round, error) {
	sqlStmt := `SELECT id FROM round WHERE user_id = ?`

	var roundIDs []int
	err := r.db.Select(&roundIDs, sqlStmt, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get round IDs: %w", err)
	}

	rounds := make([]*models.Round, 0, len(roundIDs))
	for _, id := range roundIDs {
		round, err := models.RoundById(r.db, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get round by ID: %w", err)
		}
		rounds = append(rounds, round)
	}

	return rounds, nil
}

func (r *repository) SaveRoundStats(roundStats *models.RoundStats) error {
	// Check if there is already some round stats to update.
	sqlStmt := `SELECT id FROM round_stats WHERE round_id = ?`

	var id int
	err := r.db.Get(&id, sqlStmt, roundStats.RoundId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			id = 0
		default:
			return fmt.Errorf("failed to get round stats ID: %w", err)
		}
	}

	if id == 0 {
		return roundStats.Insert(r.db)
	}

	roundStats.Id = id
	err = roundStats.Update(r.db)
	if err != nil && !errors.Is(err, models.ErrNoAffectedRows) {
		return fmt.Errorf("failed to update round stats: %w", err)
	}

	return nil
}
