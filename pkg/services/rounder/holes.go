package rounder

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/logging"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
	repo "github.com/Jacobbrewer1/golf-stats-tracker/pkg/repositories/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils"
	uhttp "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/http"
)

func (s *service) GetRoundHoles(w http.ResponseWriter, r *http.Request, roundId api.PathRoundId) {
	// Get the round by the ID.
	round, err := s.r.GetRoundById(int(roundId))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			uhttp.SendMessageWithStatus(w, http.StatusNotFound, "round not found")
			return
		default:
			slog.Error("Error getting round", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting round", err)
			return
		}
	}

	// Get the holes for the round.
	holes, err := s.r.GetRoundHoles(round.Id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNoHolesFound):
			holes = make([]*models.Hole, 0)
		default:
			slog.Error("Error getting holes", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting holes", err)
			return
		}
	}

	// Map the holes to the API model.
	respHoles := make([]*api.Hole, len(holes))
	for i, hole := range holes {
		respHoles[i] = modelHoleAsApiRoundHole(hole)
	}

	err = uhttp.Encode(w, http.StatusOK, respHoles)
	if err != nil {
		slog.Error("Error encoding holes", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func modelHoleAsApiRoundHole(hole *models.Hole) *api.Hole {
	return &api.Hole{
		Id:          utils.Ptr(int64(hole.Id)),
		Meters:      utils.Ptr(int64(hole.DistanceMeters)),
		Number:      utils.Ptr(int64(hole.Number)),
		Par:         utils.Ptr(int64(hole.Par)),
		StrokeIndex: utils.Ptr(int64(hole.Stroke)),
		Yardage:     utils.Ptr(int64(hole.DistanceYards)),
	}
}
