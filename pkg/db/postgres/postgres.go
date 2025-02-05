package postgres

import (
	"errors"

	"github.com/lib/pq"
	"github.com/omniful/go_commons/db/sql/postgres"
)

type Db struct {
	*postgres.DbCluster
}

var dbInstance *Db

func GetCluster() *Db {
	return dbInstance
}

func SetCluster(cluster *postgres.DbCluster) {
	dbInstance = &Db{cluster}
}

func IsViolatesUniqueConstraint(err error) bool {
	var pqError *pq.Error
	ok := errors.As(err, &pqError)
	if !ok || pqError == nil {
		return false
	}

	return pqError.Code == "23505"
}
