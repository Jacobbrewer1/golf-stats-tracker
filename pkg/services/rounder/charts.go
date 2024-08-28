package rounder

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/logging"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
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
			err = uhttp.Encode(w, http.StatusOK, make([]*api.ChartDataPoint, 0))
			if err != nil {
				slog.Error("Error encoding line chart data", slog.String(logging.KeyError, err.Error()))
			}
			return
		default:
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting line chart data", err)
			return
		}
	}

	// If both the from_date and since parameters are set, return an error.
	if params.FromDate != nil && params.Since != nil {
		err := uhttp.Encode(w, http.StatusBadRequest, "cannot use both from_date and since parameters")
		if err != nil {
			slog.Error("Error encoding line chart data", slog.String(logging.KeyError, err.Error()))
		}
		return
	}

	// If the "from_date" parameter is set, filter the data.
	if params.FromDate != nil {
		lineChartData = filterRoundStatsByDate(lineChartData, params.FromDate.Time)
	}

	// If the "since" parameter is set, filter the data.
	if params.Since != nil {
		duration, err := time.ParseDuration(*params.Since)
		if err != nil {
			uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "invalid duration", err)
			return
		}

		// Calculate the date that is "since" the current date.
		fromDate := time.Now().Add(-duration)
		lineChartData = filterRoundStatsByDate(lineChartData, fromDate)
	}

	data := make(map[string]float64)
	teeTimeMap := make(map[string]time.Time)

	// Fill up the map with the requested data.
	for _, d := range lineChartData {
		xVal := fmt.Sprintf("%s - %s", d.Course.Name, d.Round.TeeTime.Format(time.DateOnly))
		if _, ok := teeTimeMap[xVal]; !ok {
			teeTimeMap[xVal] = d.Round.TeeTime
		}
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
	resp := make([]*api.ChartDataPoint, 0)
	for key, value := range data {
		resp = append(resp, &api.ChartDataPoint{
			X: utils.Ptr(key),
			Y: utils.Ptr(float32(value)),
		})
	}

	sort.Slice(resp, func(i, j int) bool {
		return teeTimeMap[*resp[i].X].Before(teeTimeMap[*resp[j].X])
	})

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("Error encoding line chart data", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func filterRoundStatsByDate(data []*repo.RoundWithStats, fromDate time.Time) []*repo.RoundWithStats {
	filteredData := make([]*repo.RoundWithStats, 0)
	for _, d := range data {
		if d.Round.TeeTime.After(fromDate) {
			filteredData = append(filteredData, d)
		}
	}
	return filteredData
}

func (s *service) GetPieChartAverages(w http.ResponseWriter, r *http.Request, params api.GetPieChartAveragesParams) {
	userId := uhttp.UserIdFromContext(r.Context())

	// Get the pie chart data.
	userRounds, err := s.r.GetUserHitStats(userId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNoStatsFound):
			err = uhttp.Encode(w, http.StatusOK, make([]*api.ChartDataPoint, 0))
			if err != nil {
				slog.Error("Error encoding pie chart data", slog.String(logging.KeyError, err.Error()))
			}
			return
		default:
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting pie chart data", err)
			return
		}
	}

	data := make(map[string]int)

	hitType := ""
	switch params.AverageType {
	case api.AverageType_fairway_hit:
		hitType = models.RoundHitStatsTypeFAIRWAY
	case api.AverageType_green_hit:
		hitType = models.RoundHitStatsTypeGREEN
	}

	// Fill up the map with the requested data.
	for _, d := range userRounds {
		if string(d.Type) != hitType {
			continue
		}

		switch params.AverageType {
		case api.AverageType_fairway_hit:
			if d.Miss == models.HoleStatsFairwayHitNOTAPPLICABLE {
				continue
			}
			data[d.Miss] += d.Count
		case api.AverageType_green_hit:
			data[d.Miss] += d.Count
		}
	}

	// Create the response.
	resp := make([]*api.ChartDataPoint, 0)
	for key, value := range data {
		resp = append(resp, &api.ChartDataPoint{
			X: utils.Ptr(key),
			Y: utils.Ptr(float32(value)),
		})
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("Error encoding pie chart data", slog.String(logging.KeyError, err.Error()))
		return
	}
}
