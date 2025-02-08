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

	res, err := pd.conn.Query(ctx, "SELECT EXISTS(SELECT 1 FROM usersrecord WHERE email=$1)", email)
	if err != nil {
		return false
	}

	var exists bool
	res.Scan(&exists)

	return exists

}

func (pd *PostgresDriver) InsertUser(ctx context.Context, username string, email string) error {

	_, err := pd.conn.Exec(ctx, "INSERT INTO usersrecord(username, email) VALUES ($1, $2)", username, email)
	return err

}

func (pd *PostgresDriver) RemoveUser(ctx context.Context, email string) error {

	_, err := pd.conn.Exec(ctx, "DELETE FROM usersrecord WHERE email=$1", email)
	return err

}

func (pd *PostgresDriver) UpdateUsername(ctx context.Context, email string, newUsername string) error {

	_, err := pd.conn.Exec(ctx, "UPDATE usersrecord SET username=$1 WHERE email=$2", newUsername, email)
	return err

}
