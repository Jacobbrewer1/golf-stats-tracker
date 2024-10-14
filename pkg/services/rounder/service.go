package rounder

import (
	api "github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/rounder"
	repo "github.com/Jacobbrewer1/golf-stats-tracker/pkg/repositories/rounder"
	"github.com/Jacobbrewer1/vaulty"
	"github.com/spf13/viper"
)

type service struct {
	r            repo.Repository
	vc           vaulty.Client
	golfDataHost string
	vip          *viper.Viper
}

// NewService creates a new service.
func NewService(r repo.Repository, vc vaulty.Client, golfDataHost string, vip *viper.Viper) api.ServerInterface {
	return &service{
		r:            r,
		vc:           vc,
		golfDataHost: golfDataHost,
		vip:          vip,
	}
}
