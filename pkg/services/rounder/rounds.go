package rounder

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/logging"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
	repo "github.com/Jacobbrewer1/golf-stats-tracker/pkg/repositories/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils"
	uhttp "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/http"
	usql "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/sql"
)

func (s *service) CreateRound(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		uhttp.SendMessageWithStatus(w, http.StatusBadRequest, "request body required")
		return
	}

	rnd := new(api.RoundCreate)
	err := uhttp.DecodeJSONBody(r, rnd)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	if rnd.CourseId == nil {
		uhttp.SendMessageWithStatus(w, http.StatusBadRequest, "course_id is required")
		return
	}

	userId := uhttp.UserIdFromContext(r.Context())
	if userId <= 0 {
		slog.Debug("user_id not found in context")
		uhttp.SendMessageWithStatus(w, http.StatusUnauthorized, "user_id not found in context")
		return
	}

	mdl, err := s.roundAsModel(rnd, userId)
	if err != nil {
		uhttp.SendMessageWithStatus(w, http.StatusBadRequest, "error mapping round to model", err)
		return
	}

	err = s.r.CreateRound(mdl)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error creating round", err)
		return
	}

	err = s.importCourse(r.Context(), mdl.Id, int(*rnd.CourseId), int(*rnd.MarkerId))
	if err != nil {
		slog.Error("error importing course", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error importing course", err)
		return
	}

	respRound, err := s.roundById(mdl.Id)
	err = uhttp.Encode(w, http.StatusCreated, respRound)
	if err != nil {
		slog.Error("error encoding response", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) roundById(id int) (*api.Round, error) {
	r, err := s.r.GetRoundDetailsByRoundId(id)
	if err != nil {
		return nil, fmt.Errorf("error getting round by id: %w", err)
	}

	return s.roundAsApiRound(r), nil
}

func (s *service) roundAsApiRound(r *repo.RoundDetails) *api.Round {
	return &api.Round{
		CourseName: utils.Ptr(r.Course.Name),
		Id:         utils.Ptr(int64(r.Round.Id)),
		Marker:     utils.Ptr(r.CourseDetails.Marker.String),
		TeeTime:    utils.Ptr(r.Round.TeeTime),
	}
}

// importCourse imports a course from the golf data service.
func (s *service) importCourse(ctx context.Context, roundId int, courseId int, markerId int) error {
	if roundId == 0 {
		return errors.New("round_id is required")
	} else if courseId == 0 {
		return errors.New("course_id is required")
	} else if markerId == 0 {
		return errors.New("marker_id is required")
	}

	course, err := s.getDataCourse(ctx, courseId)
	if err != nil {
		return fmt.Errorf("failed to get course data: %w", err)
	}

	if course.Details == nil {
		return errors.New("course details are required")
	}

	details := new(models.CourseDetails)
	holes := make([]*models.Hole, 0)
	found := false
	for _, d := range *course.Details {
		if *d.Id != int64(markerId) {
			continue
		}

		courseDetails, err := s.detailsAsModel(&d)
		if err != nil {
			return fmt.Errorf("error mapping course details to model: %w", err)
		}

		if d.Holes == nil {
			return errors.New("holes are required")
		}

		for _, h := range *d.Holes {
			hole, err := s.holeAsModel(&h)
			if err != nil {
				return fmt.Errorf("error mapping hole to model: %w", err)
			}
			holes = append(holes, hole)
		}

		details = courseDetails
		found = true

		break
	}
	if !found {
		return fmt.Errorf("course details not found for marker_id: %d", markerId)
	}

	c, err := s.courseAsModel(course)
	if err != nil {
		return fmt.Errorf("error mapping course to model: %w", err)
	}
	c.RoundId = roundId

	err = s.r.CreateCourse(c)
	if err != nil {
		return fmt.Errorf("error creating course: %w", err)
	}

	details.CourseId = c.Id
	err = s.r.CreateCourseDetails(details)
	if err != nil {
		return fmt.Errorf("error creating course details: %w", err)
	}

	for _, h := range holes {
		h.CourseDetailsId = c.Id
		err = s.r.CreateHole(h)
		if err != nil {
			return fmt.Errorf("error creating hole: %w", err)
		}
	}

	return nil
}

func (s *service) holeAsModel(hole *api.Hole) (*models.Hole, error) {
	h := new(models.Hole)

	if hole.Number == nil {
		return nil, errors.New("number is required")
	}
	h.Number = int(*hole.Number)

	if hole.Par == nil {
		return nil, errors.New("par is required")
	}
	h.Par = int(*hole.Par)

	if hole.StrokeIndex == nil {
		return nil, errors.New("stroke_index is required")
	}
	h.Stroke = int(*hole.StrokeIndex)

	if hole.Yardage == nil {
		return nil, errors.New("yardage is required")
	}
	h.DistanceYards = int(*hole.Yardage)

	if hole.Meters == nil {
		return nil, errors.New("meters is required")
	}
	h.DistanceMeters = int(*hole.Meters)

	return h, nil
}

func (s *service) courseAsModel(course *api.Course) (*models.Course, error) {
	c := new(models.Course)

	if course.Name == nil {
		return nil, errors.New("name is required")
	}
	c.Name = *course.Name

	return c, nil
}

func (s *service) detailsAsModel(details *api.CourseDetails) (*models.CourseDetails, error) {
	d := new(models.CourseDetails)

	if details.Marker == nil {
		return nil, errors.New("marker is required")
	}
	d.Marker = *usql.NewNullString(*details.Marker)

	if details.Slope == nil {
		return nil, errors.New("slope is required")
	}
	d.Slope = int(*details.Slope)

	if details.Rating == nil {
		return nil, errors.New("rating is required")
	}
	d.CourseRating = *details.Rating

	if details.ParFrontNine == nil {
		return nil, errors.New("front_nine_par is required")
	}
	d.FrontNinePar = int(*details.ParFrontNine)

	if details.ParBackNine == nil {
		return nil, errors.New("back_nine_par is required")
	}
	d.BackNinePar = int(*details.ParBackNine)

	if details.ParTotal == nil {
		return nil, errors.New("total_par is required")
	}
	d.TotalPar = int(*details.ParTotal)

	if details.YardageFrontNine == nil {
		return nil, errors.New("front_nine_yards is required")
	}
	d.FrontNineYards = int(*details.YardageFrontNine)

	if details.YardageBackNine == nil {
		return nil, errors.New("back_nine_yards is required")
	}
	d.BackNineYards = int(*details.YardageBackNine)

	if details.YardageTotal == nil {
		return nil, errors.New("total_yards is required")
	}
	d.TotalYards = int(*details.YardageTotal)

	if details.MetersFrontNine == nil {
		return nil, errors.New("front_nine_meters is required")
	}
	d.FrontNineMeters = int(*details.MetersFrontNine)

	if details.MetersBackNine == nil {
		return nil, errors.New("back_nine_meters is required")
	}
	d.BackNineMeters = int(*details.MetersBackNine)

	if details.MetersTotal == nil {
		return nil, errors.New("total_meters is required")
	}
	d.TotalMeters = int(*details.MetersTotal)

	return d, nil
}

func (s *service) getDataCourse(ctx context.Context, courseId int) (*api.Course, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/courses/%d", s.golfDataHost, courseId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			slog.Error("failed to close response body", slog.String("error", err.Error()))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	course := new(api.Course)
	err = uhttp.DecodeJSON(resp.Body, course)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return course, nil
}

func (s *service) roundAsModel(rnd *api.RoundCreate, userId int) (*models.Round, error) {
	r := new(models.Round)

	if userId == 0 {
		return nil, errors.New("user_id is required")
	}
	r.UserId = userId

	if rnd.TeeTime == nil {
		return nil, errors.New("tee_time is required")
	}
	r.TeeTime = *rnd.TeeTime

	return r, nil
}

func (s *service) GetRounds(w http.ResponseWriter, r *http.Request) {
	userId := uhttp.UserIdFromContext(r.Context())
	if userId <= 0 {
		slog.Debug("user_id not found in context")
		uhttp.SendMessageWithStatus(w, http.StatusUnauthorized, "user_id not found in context")
		return
	}

	rounds, err := s.r.GetRoundsByUserId(userId)
	if err != nil {
		slog.Error("error getting rounds", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting rounds", err)
		return
	}

	respRounds := make([]*api.Round, 0)
	for _, r := range rounds {
		respRound, err := s.roundById(r.Id)
		if err != nil {
			slog.Error("error getting round by id", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting round by id", err)
			return
		}

		respRounds = append(respRounds, respRound)
	}

	err = uhttp.Encode(w, http.StatusOK, respRounds)
	if err != nil {
		slog.Error("error encoding response", slog.String(logging.KeyError, err.Error()))
		return
	}
}
