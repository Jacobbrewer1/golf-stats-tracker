package rounder

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/logging"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
	repo "github.com/Jacobbrewer1/golf-stats-tracker/pkg/repositories/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils"
	usql "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/sql"
	"github.com/Jacobbrewer1/uhttp"
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
			holes = &repo.PaginationResponse[models.Hole]{
				Items: make([]*models.Hole, 0),
				Total: 0,
			}
		default:
			slog.Error("Error getting holes", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting holes", err)
			return
		}
	}

	// Map the holes to the API model.
	respHoles := make([]api.Hole, len(holes.Items))
	for i, hole := range holes.Items {
		respHoles[i] = *modelHoleAsApiRoundHole(hole)
	}

	resp := &api.HolesResponse{
		Holes: &respHoles,
		Total: utils.Ptr(holes.Total),
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
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
		case errors.Is(err, repo.ErrNoStatsFound):
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
		case errors.Is(err, repo.ErrNoStatsFound):
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
	opts := cmpopts.IgnoreFields(models.HoleStats{}, "Id", "HoleId")
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

	go func() {
		csErr := s.calculateStats(round.UserId, round.Id)
		if csErr != nil {
			slog.Error("Error calculating stats", slog.String(logging.KeyError, csErr.Error()))
		}
	}()

	err = uhttp.Encode(w, http.StatusOK, modelHoleStatsAsApiHoleStats(newStats))
	if err != nil {
		slog.Error("Error encoding hole stats", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) calculateStats(userId int, roundId int) error {
	// Get the line chart data.
	roundData, err := s.r.GetStatsByRoundId(userId, roundId)
	if err != nil {
		return fmt.Errorf("error getting line chart data: %w", err)
	}

	// Calculate the average score per round.
	totalFairwayCount := 0
	totalFairwayHit := 0
	totalGreenHit := 0
	totalPutts := 0
	penalties := 0
	totalScorePar3 := 0
	totalPar3 := 0
	totalScorePar4 := 0
	totalPar4 := 0
	totalScorePar5 := 0
	totalPar5 := 0

	for _, data := range roundData.Items {
		penalties += data.Stats.Penalties
		totalPutts += data.Stats.Putts
		if string(data.Stats.FairwayHit) != models.HoleStatsFairwayHitNOTAPPLICABLE {
			totalFairwayCount += 1
		}
		if string(data.Stats.FairwayHit) == models.HoleStatsFairwayHitHIT {
			totalFairwayHit += 1
		}
		if string(data.Stats.GreenHit) == models.HoleStatsGreenHitHIT {
			totalGreenHit += 1
		}

		switch data.Hole.Par {
		case 3:
			totalScorePar3 += data.Stats.Score
			totalPar3 += 1
		case 4:
			totalScorePar4 += data.Stats.Score
			totalPar4 += 1
		case 5:
			totalScorePar5 += data.Stats.Score
			totalPar5 += 1
		}
	}

	// Calculate the averages
	averagePutts := float64(totalPutts) / float64(len(roundData.Items))
	averageFairwayHit := (float64(totalFairwayHit) / float64(totalFairwayCount)) * 100
	averageGreenHit := (float64(totalGreenHit) / float64(len(roundData.Items))) * 100
	avgPar3 := float64(totalScorePar3) / float64(totalPar3)
	avgPar4 := float64(totalScorePar4) / float64(totalPar4)
	avgPar5 := float64(totalScorePar5) / float64(totalPar5)

	m := &models.RoundStats{
		RoundId:        roundId,
		AvgFairwaysHit: averageFairwayHit,
		AvgGreensHit:   averageGreenHit,
		AvgPutts:       averagePutts,
		Penalties:      penalties,
		AvgPar3:        avgPar3,
		AvgPar4:        avgPar4,
		AvgPar5:        avgPar5,
	}

	err = s.r.SaveRoundStats(m)
	if err != nil {
		return fmt.Errorf("error saving round stats: %w", err)
	}

	err = s.calculatePieStats(roundId)
	if err != nil {
		return fmt.Errorf("error calculating pie stats: %w", err)
	}

	return nil
}

func (s *service) calculatePieStats(roundId int) error {
	holes, err := s.r.GetRoundHoles(roundId)
	if err != nil {
		return fmt.Errorf("error getting holes: %w", err)
	}

	roundStats, err := s.r.GetRoundStatsByRoundId(roundId)
	if err != nil {
		return fmt.Errorf("error getting round stats: %w", err)
	}

	fairwayData := make(map[string]int)
	greenData := make(map[string]int)

	for _, h := range holes.Items {
		stats, err := s.r.GetHoleStatsByHoleId(h.Id)
		if err != nil {
			return fmt.Errorf("error getting hole stats: %w", err)
		}

		fairwayData[string(stats.FairwayHit)] += 1
		greenData[string(stats.GreenHit)] += 1
	}

	currentHitStats, err := s.r.GetRoundHitStatsByRoundStatsId(roundId)
	if err != nil {
		return fmt.Errorf("error getting current hit stats: %w", err)
	}

	fairwayHitStatIds := make(map[string]int)
	greenHitStatIds := make(map[string]int)

	for _, hitStat := range currentHitStats.Items {
		switch string(hitStat.Type) {
		case models.RoundHitStatsTypeFAIRWAY:
			fairwayHitStatIds[hitStat.Miss] = hitStat.Id
		case models.RoundHitStatsTypeGREEN:
			greenHitStatIds[hitStat.Miss] = hitStat.Id
		}
	}

	wg := new(sync.WaitGroup)
	fairwayModel := make([]*models.RoundHitStats, 0)
	wg.Add(1)
	go func() {
		for k, v := range fairwayData {
			id := 0
			if val, ok := fairwayHitStatIds[k]; ok {
				id = val
			}

			fairwayModel = append(fairwayModel, &models.RoundHitStats{
				Id:           id,
				Type:         usql.NewEnum(models.RoundHitStatsTypeFAIRWAY),
				Miss:         k,
				Count:        v,
				RoundStatsId: roundStats.Id,
			})
		}
		wg.Done()
	}()

	greenModel := make([]*models.RoundHitStats, 0)
	wg.Add(1)
	go func() {
		for k, v := range greenData {
			id := 0
			if val, ok := greenHitStatIds[k]; ok {
				id = val
			}

			greenModel = append(greenModel, &models.RoundHitStats{
				Id:           id,
				Type:         usql.NewEnum(models.RoundHitStatsTypeGREEN),
				Miss:         k,
				Count:        v,
				RoundStatsId: roundStats.Id,
			})
		}
		wg.Done()
	}()

	wg.Wait()

	err = s.r.SaveRoundHitStats(fairwayModel...)
	if err != nil {
		return fmt.Errorf("error saving fairway stats: %w", err)
	}

	err = s.r.SaveRoundHitStats(greenModel...)
	if err != nil {
		return fmt.Errorf("error saving green stats: %w", err)
	}

	return nil
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
	s.FairwayHit = usql.NewEnum(strings.ToUpper(*stats.FairwayHit))

	if stats.GreenHit == nil {
		return nil, errors.New("green_hit is required")
	} else if *stats.GreenHit == api.HitInRegulation_not_applicable {
		return nil, errors.New("green_hit cannot be NOT_APPLICABLE")
	}
	s.GreenHit = usql.NewEnum(strings.ToUpper(*stats.GreenHit))

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
