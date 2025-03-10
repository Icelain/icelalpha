package database

import (
	"context"
	"sync"
)

func Sync(postgresDriver *PostgresDriver, cache sync.Map) {

	cache.Range(func(key, value any) bool {

	})

}
