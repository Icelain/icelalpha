package database

import (
	"context"
	"sync"
)

func Sync(postgresDriver *PostgresDriver, cache *sync.Map) error {

	var err error

	cache.Range(func(key, value any) bool {

		if err = postgresDriver.UpdateUserCredits(context.Background(), key.(string), value.(uint)); err != nil {

			return false

		}

		return true

	})

	return err

}
