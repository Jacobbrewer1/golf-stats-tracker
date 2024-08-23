package rounder

import (
	"errors"
	"fmt"
	"net/http"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	repo "github.com/Jacobbrewer1/golf-stats-tracker/pkg/repositories/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils"
	uhttp "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/http"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/vault"
)

type authz struct {
	next api.ServerInterface
	db   repo.Repository
	vc   vault.Client
}

func (a *authz) CreateRound(w http.ResponseWriter, r *http.Request) {
	r, err := a.WithAuthorization(r)
	if err != nil {
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
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.GetNewRoundCourses(w, r, params)
}

func (a *authz) GetNewRoundMarker(w http.ResponseWriter, r *http.Request, courseId api.PathCourseId) {
	r, err := a.WithAuthorization(r)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusUnauthorized, "failed to authorize request", err)
		return
	}

	a.next.GetNewRoundMarker(w, r, courseId)
}

func (a *authz) CreateUser(w http.ResponseWriter, r *http.Request) {
	a.next.CreateUser(w, r)
}

func NewAuthz(next api.ServerInterface, db repo.Repository, vc vault.Client) api.ServerInterface {
	return &authz{
		next: next,
		db:   db,
		vc:   vc,
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

	unvaultedPassword, err := a.vc.TransitDecrypt(r.Context(), user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	if !utils.ComparePassword(unvaultedPassword, password) {
		return nil, errors.New("invalid password")
	}

	r = r.WithContext(uhttp.UserIdToContext(r.Context(), user.Id))

	return r, nil
}
