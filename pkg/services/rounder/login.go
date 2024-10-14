package rounder

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/logging"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils"
	"github.com/Jacobbrewer1/uhttp"
	"github.com/Jacobbrewer1/vaulty"
)

func (s *service) Login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		slog.Debug("basic auth not provided")
		uhttp.SendMessageWithStatus(w, http.StatusUnauthorized, "basic auth required")
		return
	}

	// Get the user from the database
	user, err := s.r.UserByUsername(strings.ToLower(username))
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting user", err)
		return
	}

	// Check the password
	err = s.checkPassword(r.Context(), password, user.Password)
	if err != nil {
		slog.Debug("error checking password", slog.String(logging.KeyError, err.Error()))
		uhttp.SendMessageWithStatus(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	mockRequest, err := http.NewRequest(http.MethodPost, "http://localhost:8200/v1/auth/token/lookup-self", nil)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error creating request", err)
		return
	}
	mockRequest.SetBasicAuth(user.Username, password)

	t := new(api.Token)
	// Trim the Basic prefix from the Authorization header
	t.Token = utils.Ptr(strings.TrimPrefix(mockRequest.Header.Get("Authorization"), "Basic "))

	err = uhttp.Encode(w, http.StatusOK, t)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error encoding response", err)
		return
	}
}

func (s *service) checkPassword(ctx context.Context, password, hashedPassword string) error {
	unhashedPassword, err := s.vc.Path(
		s.vip.GetString("vault.transit.key"),
		vaulty.WithPrefix(s.vip.GetString("vault.transit.name")),
	).TransitDecrypt(ctx, hashedPassword)
	if err != nil {
		return fmt.Errorf("error decrypting password: %w", err)
	}

	if !utils.ComparePassword(unhashedPassword, password) {
		return fmt.Errorf("invalid password")
	}
	return nil
}
