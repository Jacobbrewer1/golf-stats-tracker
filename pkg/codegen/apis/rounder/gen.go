package rounder

//go:generate oapi-codegen -generate types -package rounder -templates ../../templates -o types.go -import-mapping=../common/common.yaml:github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/common ./routes.yaml
//go:generate oapi-codegen -generate gorilla -package rounder -templates ../../templates -o server.go -import-mapping=../common/common.yaml:github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/common ./routes.yaml
//go:generate oapi-codegen -generate client -package rounder -templates ../../templates -o client.go -import-mapping=../common/common.yaml:github.com/Jacobbrewer1/golf-stats-tracker/pkg/codegen/apis/common ./routes.yaml
