//go:build wireinject
// +build wireinject

package hubs

import (
	"context"

	"github.com/google/wire"
	"github.com/omniful/go_commons/db/sql/postgres"
	oredis "github.com/omniful/go_commons/redis"
)

func Wire(ctx context.Context, db *postgres.DbCluster, redis *oredis.Client, nameSpace string) (*Controller, error) {
	panic(wire.Build(ProviderSet))
}
