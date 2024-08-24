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
	usql "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/sql"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	} else if round.UserId != uhttp.UserIdFromContext(r.Context()) {
		uhttp.SendMessageWithStatus(w, http.StatusForbidden, "round not found")
		return
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

func (s *service) GetHoleStats(w http.ResponseWriter, r *http.Request, roundId api.PathRoundId, holeId api.PathHoleId) {
	// Get the round by the ID.
	round, err := s.r.GetRoundById(int(roundId))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Debug("round not found", slog.Int64("round_id", roundId))
			uhttp.SendMessageWithStatus(w, http.StatusNotFound, "round not found")
			return
		default:
			slog.Error("Error getting round", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting round", err)
			return
		}
	} else if round.UserId != uhttp.UserIdFromContext(r.Context()) {
		uhttp.SendMessageWithStatus(w, http.StatusForbidden, "round not found")
		return
	}

	// Get the hole by the ID.
	hole, err := s.r.GetHoleById(int(holeId))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Debug("hole not found", slog.Int("round_id", round.Id), slog.Int64("hole_id", holeId))
			uhttp.SendMessageWithStatus(w, http.StatusNotFound, "hole not found")
			return
		default:
			slog.Error("Error getting hole", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting hole", err)
			return
		}
	}

	// Get the stats for the hole.
	stats, err := s.r.GetHoleStatsByHoleId(hole.Id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNoHoleStatsFound):
			stats = new(models.HoleStats)
			stats.HoleId = hole.Id
		default:
			slog.Error("Error getting hole stats", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting hole stats", err)
			return
		}
	}

	// Map the stats to the API model.
	respStats := modelHoleStatsAsApiHoleStats(stats)

	err = uhttp.Encode(w, http.StatusOK, respStats)
	if err != nil {
		slog.Error("Error encoding hole stats", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func modelHoleStatsAsApiHoleStats(stats *models.HoleStats) *api.HoleStats {
	s := new(api.HoleStats)
	s.Score = utils.Ptr(int64(stats.Score))
	s.FairwayHit = utils.Ptr(api.HitInRegulation(stats.FairwayHit))
	s.GreenHit = utils.Ptr(api.HitInRegulation(stats.GreenHit))
	s.Putts = utils.Ptr(int64(stats.Putts))
	s.Penalties = utils.Ptr(int64(stats.Penalties))
	s.PinLocation = utils.Ptr(stats.PinLocation)
	return s
}

func (s *service) UpdateHoleStats(w http.ResponseWriter, r *http.Request, roundId api.PathRoundId, holeId api.PathHoleId) {
	if r.Body == http.NoBody {
		uhttp.SendMessageWithStatus(w, http.StatusBadRequest, "request body required")
		return
	}

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
	} else if round.UserId != uhttp.UserIdFromContext(r.Context()) {
		uhttp.SendMessageWithStatus(w, http.StatusForbidden, "round not found")
		return
	}

	// Get the hole by the ID.
	hole, err := s.r.GetHoleById(int(holeId))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			uhttp.SendMessageWithStatus(w, http.StatusNotFound, "hole not found")
			return
		default:
			slog.Error("Error getting hole", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting hole", err)
			return
		}
	}

	// Decode the request body into the API model.
	reqStats := new(api.HoleStats)
	err = uhttp.DecodeJSON(r.Body, reqStats)
	if err != nil {
		slog.Error("Error decoding request body", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	// Get the stats for the hole.
	stats, err := s.r.GetHoleStatsByHoleId(hole.Id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNoHoleStatsFound):
			stats = new(models.HoleStats)
			stats.HoleId = hole.Id
		default:
			slog.Error("Error getting hole stats", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting hole stats", err)
			return
		}
	}

	// Map the API model to the model.
	newStats, err := apiAsModelHoleStats(reqStats)
	if err != nil {
		slog.Error("Error mapping hole stats to model", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "error mapping hole stats to model", err)
		return
	}

	// Is there any difference between the existing stats and the new stats?
	opts := cmpopts.IgnoreFields(models.HoleStats{}, "Id")
	if diff := cmp.Diff(stats, newStats, opts); diff == "" {
		slog.Debug("hole stats are the same", slog.Int("round_id", round.Id), slog.Int("hole_id", hole.Id))
	} else {
		slog.Debug("hole stats are different", slog.String("diff", diff))

		if stats.Id != 0 {
			newStats.Id = stats.Id
		}
		newStats.HoleId = stats.HoleId

		err = s.r.SaveHoleStats(newStats)
		if err != nil {
			slog.Error("Error saving hole stats", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error saving hole stats", err)
			return
		}
	}

	err = uhttp.Encode(w, http.StatusOK, modelHoleStatsAsApiHoleStats(newStats))
	if err != nil {
		slog.Error("Error encoding hole stats", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func apiAsModelHoleStats(stats *api.HoleStats) (*models.HoleStats, error) {
	s := new(models.HoleStats)
	if stats.Score == nil {
		return nil, errors.New("score is required")
	}
	s.Score = int(*stats.Score)

	if stats.FairwayHit == nil {
		return nil, errors.New("fairway_hit is required")
	}
	s.FairwayHit = usql.NewEnum(*stats.FairwayHit)

	if stats.GreenHit == nil {
		return nil, errors.New("green_hit is required")
	} else if *stats.GreenHit == api.HitInRegulation_NOT_APPLICABLE {
		return nil, errors.New("green_hit cannot be NOT_APPLICABLE")
	}
	s.GreenHit = usql.NewEnum(*stats.GreenHit)

	if stats.Putts == nil {
		return nil, errors.New("putts is required")
	}
	s.Putts = int(*stats.Putts)

	if stats.Penalties == nil {
		return nil, errors.New("penalties is required")
	}
	s.Penalties = int(*stats.Penalties)

	if stats.PinLocation == nil {
		return nil, errors.New("pin_location is required")
	}
	s.PinLocation = *stats.PinLocation

	return s, nil
}
