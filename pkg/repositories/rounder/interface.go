package rounder

import "github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"

type Repository interface {
	// CreateUser creates a new user.
	CreateUser(user *models.User) error

	// UserByUsername returns the user with the given username.
	UserByUsername(username string) (*models.User, error)

	// CreateRound creates a new round.
	CreateRound(round *models.Round) error

	// CreateCourse creates a new course.
	CreateCourse(course *models.Course) error

	// CreateCourseDetails creates a new course details.
	CreateCourseDetails(courseDetails *models.CourseDetails) error

	// CreateHole creates a new hole.
	CreateHole(hole *models.Hole) error

	// GetRoundById gets a round by its ID.
	GetRoundById(id int) (*models.Round, error)

	// GetRoundDetailsByRoundId gets the details for a round.
	GetRoundDetailsByRoundId(roundId int) (*RoundDetails, error)

	// GetRoundsByUserId gets the rounds for a user.
	GetRoundsByUserId(userId int) ([]*models.Round, error)

	// GetRoundHoles gets the holes for a round.
	GetRoundHoles(roundId int) ([]*models.Hole, error)

	// GetHoleById gets a hole by its ID.
	GetHoleById(id int) (*models.Hole, error)

	// GetHoleStatsByHoleId gets the stats for a hole.
	GetHoleStatsByHoleId(holeId int) (*models.HoleStats, error)

	// SaveHoleStats saves the stats for a hole.
	SaveHoleStats(holeStats *models.HoleStats) error

	// GetAllStatsForPar gets all stats for a user.
	GetAllStatsForPar(userId int, par int64) ([]*HoleWithStats, error)

	// CountHolesByRoundAndPar counts the number of holes for a round and par.
	CountHolesByRoundAndPar(roundId int, par int64) (int, error)

	// GetStatsByRoundId gets the stats for a user.
	GetStatsByRoundId(userId int, roundId int) ([]*HoleWithStats, error)

	// CountHolesByRoundId counts the number of holes for a round.
	CountHolesByRoundId(roundId int) (int, error)

	// SaveRoundStats saves the stats for a round.
	SaveRoundStats(roundStats *models.RoundStats) error

	// GetStatsByUserId gets all stats for a user.
	GetStatsByUserId(userId int) ([]*RoundWithStats, error)
}

type HoleWithStats struct {
	Round *models.Round
	Hole  *models.Hole
	Stats *models.HoleStats
}

type RoundWithStats struct {
	Round  *models.Round
	Course *models.Course
	Stats  *models.RoundStats
}

type RoundDetails struct {
	Round         *models.Round
	Course        *models.Course
	CourseDetails *models.CourseDetails
	Holes         []*models.Hole
}
