package rounder

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/logging"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils"
	uhttp "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/http"
)

func (s *service) GetLineChartScoreAverage(w http.ResponseWriter, r *http.Request, params api.GetLineChartScoreAverageParams) {
	userId := uhttp.UserIdFromContext(r.Context())

	// Get the line chart data.
	lineChartData, err := s.r.GetAllStatsForPar(userId, params.Par)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting line chart data", err)
		return
	}

	// Calculate the average score per round.
	averageScores := make(map[int]float64)
	for _, data := range lineChartData {
		if data.Hole.Par != int(params.Par) {
			continue
		}

		averageScores[data.Round.Id] += float64(data.Stats.Score)
	}

	// Find how many holes of that par there are.
	for key, value := range averageScores {
		count, err := s.r.CountHolesByRoundAndPar(key, params.Par)
		if err != nil {
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting hole count", err)
			return
		}

		averageScores[key] = value / float64(count)
	}

	sort.Slice(lineChartData, func(i, j int) bool {
		return lineChartData[i].Round.TeeTime.Before(lineChartData[j].Round.TeeTime)
	})

	// Create the response.
	resp := make([]*api.LineDataPoint, 0)
	for key, value := range averageScores {
		round, err := s.r.GetRoundDetailsByRoundId(key)
		if err != nil {
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting round", err)
			return
		}

		roundedValue := utils.Round(value, 2)

		resp = append(resp, &api.LineDataPoint{
			X: utils.Ptr(fmt.Sprintf("%s - %s", round.Course.Name, round.Round.TeeTime.Format(time.DateOnly))),
			Y: utils.Ptr(float32(roundedValue)),
		})
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("Error encoding line chart data", slog.String(logging.KeyError, err.Error()))
		return
	}
}
