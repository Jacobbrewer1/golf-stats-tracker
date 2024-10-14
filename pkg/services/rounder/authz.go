package rounder

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/logging"
	repo "github.com/Jacobbrewer1/golf-stats-tracker/pkg/repositories/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils"
	"github.com/Jacobbrewer1/uhttp"
	"github.com/Jacobbrewer1/vaulty"
	"github.com/spf13/viper"
)

type authz struct {
	next api.ServerInterface
	db   repo.Repository
	vc   vaulty.Client
	vip  *viper.Viper
}

func (a *authz) GetPieChartAverages(w http.ResponseWriter, r *http.Request, params api.GetPieChartAveragesParams) {
	r, err := a.WithAuthorization(r)
	if err != nil {
		slog.Debug("failed to authorize request", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.GetPieChartAverages(w, r, params)
}

func (a *authz) GetLineChartAverages(w http.ResponseWriter, r *http.Request, params api.GetLineChartAveragesParams) {
	r, err := a.WithAuthorization(r)
	if err != nil {
		slog.Debug("failed to authorize request", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.GetLineChartAverages(w, r, params)
}

func (a *authz) UpdateHoleStats(w http.ResponseWriter, r *http.Request, roundId api.PathRoundId, holeId api.PathHoleId) {
	r, err := a.WithAuthorization(r)
	if err != nil {
		slog.Debug("failed to authorize request", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.UpdateHoleStats(w, r, roundId, holeId)
}

func (a *authz) GetHoleStats(w http.ResponseWriter, r *http.Request, roundId api.PathRoundId, holeId api.PathHoleId) {
	r, err := a.WithAuthorization(r)
	if err != nil {
		slog.Debug("failed to authorize request", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.GetHoleStats(w, r, roundId, holeId)
}

func (a *authz) GetRoundHoles(w http.ResponseWriter, r *http.Request, roundId api.PathRoundId) {
	r, err := a.WithAuthorization(r)
	if err != nil {
		slog.Debug("failed to authorize request", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.GetRoundHoles(w, r, roundId)
}

func (a *authz) GetRounds(w http.ResponseWriter, r *http.Request) {
	r, err := a.WithAuthorization(r)
	if err != nil {
		slog.Debug("failed to authorize request", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.GetRounds(w, r)
}

func (a *authz) CreateRound(w http.ResponseWriter, r *http.Request) {
	r, err := a.WithAuthorization(r)
	if err != nil {
		slog.Debug("failed to authorize request", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.CreateRound(w, r)
}

func (a *authz) Login(w http.ResponseWriter, r *http.Request) {
	a.next.Login(w, r)
}

func (a *authz) GetNewRoundCourses(w http.ResponseWriter, r *http.Request, params api.GetNewRoundCoursesParams) {
	r, err := a.WithAuthorization(r)
	if err != nil {
		slog.Debug("failed to authorize request", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.GetNewRoundCourses(w, r, params)
}

func (a *authz) GetNewRoundMarker(w http.ResponseWriter, r *http.Request, courseId api.PathCourseId) {
	r, err := a.WithAuthorization(r)
	if err != nil {
		slog.Debug("failed to authorize request", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.GetNewRoundMarker(w, r, courseId)
}

func (a *authz) CreateUser(w http.ResponseWriter, r *http.Request) {
	a.next.CreateUser(w, r)
}

func NewAuthz(next api.ServerInterface, db repo.Repository, vc vaulty.Client, vip *viper.Viper) api.ServerInterface {
	return &authz{
		next: next,
		db:   db,
		vc:   vc,
		vip:  vip,
	}
}

func (a *authz) WithAuthorization(r *http.Request) (*http.Request, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, errors.New("missing basic auth")
	}

	user, err := a.db.UserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	unvaultedPassword, err := a.vc.Path(
		a.vip.GetString("vault.transit.key"),
		vaulty.WithPrefix(a.vip.GetString("vault.transit.name")),
	).TransitDecrypt(r.Context(), user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	if !utils.ComparePassword(unvaultedPassword, password) {
		return nil, errors.New("invalid password")
	}

	r = r.WithContext(utils.UserIdToContext(r.Context(), user.Id))

	return r, nil
}
