package rounder

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
	uhttp "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/http"
)

func (s *service) CreateRound(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		uhttp.SendMessageWithStatus(w, http.StatusBadRequest, "request body required")
		return
	}

	rnd := new(api.Round)
	err := uhttp.DecodeJSONBody(r, rnd)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	if rnd.CourseId == nil {
		uhttp.SendMessageWithStatus(w, http.StatusBadRequest, "course_id is required")
		return
	}

	mdl, err := s.roundAsModel(rnd)
	if err != nil {
		uhttp.SendMessageWithStatus(w, http.StatusBadRequest, "error mapping round to model", err)
		return
	}

	err = s.r.CreateRound(mdl)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error creating round", err)
		return
	}
}

// importCourse imports a course from the golf data service.
func (s *service) importCourse(ctx context.Context, courseId int, markerId int) error {
	course, err := s.getDataCourse(ctx, courseId)
	if err != nil {
		return fmt.Errorf("failed to get course data: %w", err)
	}

	if course.Details == nil {
		return errors.New("course details are required")
	}

	for i := range *course.Details {
		dets := (*course.Details)[i]
		if *dets.Id != int64(markerId) {
			continue
		}
	}
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

func (s *service) roundAsModel(rnd *api.Round) (*models.Round, error) {
	r := new(models.Round)

	if rnd.MarkerId == nil {
		return nil, errors.New("marker_id is required")
	}
	r.CourseId = int(*rnd.CourseId)

	if rnd.TeeTime == nil {
		return nil, errors.New("tee_time is required")
	}
	r.TeeTime = *rnd.TeeTime

	return r, nil
}
