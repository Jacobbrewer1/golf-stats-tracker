package rounder

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/logging"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
	repo "github.com/Jacobbrewer1/golf-stats-tracker/pkg/repositories/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils"
	uhttp "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/http"
	"github.com/Jacobbrewer1/vaulty"
)

func (s *service) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		uhttp.SendMessageWithStatus(w, http.StatusBadRequest, "request body is required")
		return
	}

	user := new(api.User)
	err := uhttp.DecodeJSONBody(r, user)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}
	user.Id = nil

	// Ensure that the username is unique
	_, err = s.r.UserByUsername(strings.ToLower(*user.Username))
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrUserNotFound):
			// Continue
		default:
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error checking if user exists", err)
			return
		}
	} else {
		uhttp.SendMessageWithStatus(w, http.StatusConflict, "user already exists")
		return
	}

	u, err := s.userAsModel(*user)
	if err != nil {
		uhttp.SendHttpError(w, err)
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}
	u.Password = hashedPassword

	// Encrypt the password with vault
	encryptedPassword, err := s.vc.Path(
		s.vip.GetString("vault.transit.key"),
		vaulty.WithPrefix(s.vip.GetString("vault.transit.name")),
	).TransitEncrypt(r.Context(), u.Password)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error encrypting password", err)
		return
	}
	u.Password = vaulty.GetTransitCipherText(encryptedPassword)

	if err := s.r.CreateUser(u); err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error creating user", err)
		return
	}

	err = uhttp.Encode(w, http.StatusCreated, s.modelAsUser(u))
	if err != nil {
		slog.Error("error encoding response", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) modelAsUser(u *models.User) *api.User {
	return &api.User{
		Id:       utils.Ptr(int64(u.Id)),
		Name:     utils.Ptr(u.Name),
		Username: utils.Ptr(u.Username),
	}
}

func (s *service) userAsModel(user api.User) (*models.User, error) {
	u := new(models.User)

	if user.Username == nil {
		return nil, uhttp.NewHttpError(http.StatusBadRequest, "username is required")
	}
	u.Username = strings.ToLower(*user.Username)

	if user.Password == nil {
		return nil, uhttp.NewHttpError(http.StatusBadRequest, "password is required")
	}
	u.Password = *user.Password

	if user.Name == nil {
		return nil, uhttp.NewHttpError(http.StatusBadRequest, "name is required")
	}
	u.Name = *user.Name

	return u, nil
}
