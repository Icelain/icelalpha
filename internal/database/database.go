package database

import (
	"context"
	"icealpha/internal/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Db interface will be mocked once initial release is done

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

func (pd *PostgresDriver) GetUser(ctx context.Context, email string) (types.User, error) {

	result := types.User{}
	res, err := pd.conn.Query(ctx, "SELECT (uuid, username, email, credits) FROM usersrecord WHERE email=$1", email)

	if err != nil {

		return types.User{}, err

	}

	pgxUuid := pgtype.UUID{}
	res.Scan(&pgxUuid, &result.Username, &result.Email, &result.CreditBalance)
	result.UUID, err = uuid.FromBytes(pgxUuid.Bytes[:])

	if err != nil {

		return types.User{}, err

	}

	return result, nil

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

func (pd *PostgresDriver) UpdateUserCredits(ctx context.Context, email string) error {

	_, err := pd.conn.Exec(ctx, "UPDATE usersrecord SET credits=credits-1 WHERE email=$1", email)
	return err

}
