package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PostgresDriver struct {
	conn *pgx.Conn
}

func (pd *PostgresDriver) Close(ctx context.Context) error {

	return pd.conn.Close(ctx)

}

func CreatePostgresDriver(connectionURL string) (*PostgresDriver, error) {

	conn, err := pgx.Connect(context.Background(), connectionURL)
	if err != nil {
		return &PostgresDriver{}, err
	}

	return &PostgresDriver{

		conn: conn,
	}, nil

}

func (pd *PostgresDriver) CheckUserExists(ctx context.Context, email string) bool {

	res, err := pd.conn.Query(ctx, "select from users where email=", email)
	if err != nil {
		return false
	}

	var exists bool
	res.Scan(&exists)

	return exists

}
