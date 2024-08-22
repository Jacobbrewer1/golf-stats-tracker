package rounder

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils"
	uhttp "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/http"
)

func (s *service) Login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		uhttp.SendMessageWithStatus(w, http.StatusUnauthorized, "basic auth required")
		return
	}

	// Get the user from the database
	user, err := s.r.UserByUsername(username)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting user", err)
		return
	}

	// Check the password
	err = s.checkPassword(r.Context(), password, user.Password)
	if err != nil {
		uhttp.SendMessageWithStatus(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	tokenStr := fmt.Sprintf("%s:%s", user.Username, user.Password)
	token := base64.StdEncoding.EncodeToString([]byte(tokenStr))

	t := new(api.Token)
	t.Token = &token

	err = uhttp.Encode(w, http.StatusOK, t)
	if err != nil {
		uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error encoding response", err)
		return
	}
}

func (s *service) checkPassword(ctx context.Context, password, hashedPassword string) error {
	unhashedPassword, err := s.vc.TransitDecrypt(ctx, hashedPassword)
	if err != nil {
		return fmt.Errorf("error decrypting password: %w", err)
	}

	if !utils.ComparePassword(unhashedPassword, password) {
		return fmt.Errorf("invalid password")
	}
	return nil
}
