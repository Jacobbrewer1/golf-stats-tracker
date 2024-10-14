package rounder

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sort"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/logging"
	"github.com/Jacobbrewer1/uhttp"
)

func (s *service) GetNewRoundCourses(w http.ResponseWriter, r *http.Request, params api.GetNewRoundCoursesParams) {
	name := ""
	if params.Name != nil {
		name = *params.Name
	}

	courses, err := s.getGolfDataCourses(r.Context(), name)
	if err != nil {
		slog.Error("failed to get courses", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "failed to get courses", err)
		return
	}

	resp := &api.CoursesResponse{
		Courses: courses.Courses,
		Total:   courses.Total,
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("failed to encode response", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) getGolfDataCourses(ctx context.Context, name string) (*api.CoursesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/courses", s.golfDataHost), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if name != "" {
		params := req.URL.Query()
		params.Add("name", name)
		req.URL.RawQuery = params.Encode()
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

	courses := new(api.CoursesResponse)
	err = uhttp.DecodeJSON(resp.Body, &courses)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return courses, nil
}

func (s *service) GetNewRoundMarker(w http.ResponseWriter, r *http.Request, courseId api.PathCourseId) {
	marker, err := s.getGolfDataMarker(r.Context(), int(courseId))
	if err != nil {
		slog.Error("failed to get marker", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "failed to get marker", err)
		return
	}

	details := make([]api.CourseDetails, 0)
	details = append(details, *marker.Details...)

	for i := range details {
		details[i].Holes = nil
	}

	sort.Slice(details, func(i, j int) bool {
		if *details[i].Slope == *details[j].Slope {
			return *details[i].Rating < *details[j].Rating
		}
		return *details[i].Slope < *details[j].Slope
	})

	resp := &api.CourseDetailsResponse{
		Details: details,
		Total:   int64(len(details)),
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("failed to encode response", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) getGolfDataMarker(ctx context.Context, courseId int) (*api.Course, error) {
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

	marker := new(api.Course)
	err = uhttp.DecodeJSON(resp.Body, marker)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return marker, nil
}
