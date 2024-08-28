package rounder

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/logging"
	repo "github.com/Jacobbrewer1/golf-stats-tracker/pkg/repositories/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils"
	uhttp "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/http"
)

func (s *service) GetLineChartAverages(w http.ResponseWriter, r *http.Request, params api.GetLineChartAveragesParams) {
	userId := uhttp.UserIdFromContext(r.Context())

	// Get the line chart data.
	lineChartData, err := s.r.GetStatsByUserId(userId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNoStatsFound):
			err = uhttp.Encode(w, http.StatusOK, []*api.LineDataPoint{})
			if err != nil {
				slog.Error("Error encoding line chart data", slog.String(logging.KeyError, err.Error()))
			}
			return
		default:
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting line chart data", err)
			return
		}
	}

	data := make(map[string]float64)

	// Fill up the map with the requested data.
	for _, d := range lineChartData {
		xVal := fmt.Sprintf("%s - %s", d.Course.Name, d.Round.TeeTime.Format(time.DateOnly))
		switch params.AverageType {
		case api.AverageType_fairway_hit:
			data[xVal] += d.Stats.AvgFairwaysHit
		case api.AverageType_green_hit:
			data[xVal] += d.Stats.AvgGreensHit
		case api.AverageType_putts:
			data[xVal] += d.Stats.AvgPutts
		case api.AverageType_penalties:
			data[xVal] += float64(d.Stats.Penalties)
		case api.AverageType_par_3:
			data[xVal] += d.Stats.AvgPar3
		case api.AverageType_par_4:
			data[xVal] += d.Stats.AvgPar4
		case api.AverageType_par_5:
			data[xVal] += d.Stats.AvgPar5
		}
	}

	// Create the response.
	resp := make([]*api.LineDataPoint, 0)
	for key, value := range data {
		resp = append(resp, &api.LineDataPoint{
			X: utils.Ptr(key),
			Y: utils.Ptr(float32(value)),
		})
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("Error encoding line chart data", slog.String(logging.KeyError, err.Error()))
		return
	}
}
