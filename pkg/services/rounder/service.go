package rounder

import (
	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	repo "github.com/Jacobbrewer1/golf-stats-tracker/pkg/repositories/rounder"
	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/vault"
)

type service struct {
	r            repo.Repository
	vc           vault.Client
	golfDataHost string
}

// NewService creates a new service.
func NewService(r repo.Repository, vc vault.Client, golfDataHost string) api.ServerInterface {
	return &service{
		r:            r,
		vc:           vc,
		golfDataHost: golfDataHost,
	}
}
