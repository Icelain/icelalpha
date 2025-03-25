package database

import (
	"context"
	"icealpha/internal/types"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const MIGRATIONDIR = "./migrations"

// Try migrating if a new db is connected
func tryMigrate(connection *pgx.Conn) error {

	var globalErr error

	filepath.Walk(MIGRATIONDIR, func(fp string, info fs.FileInfo, err error) error {

		fileContent, err := os.ReadFile(fp)
		if err != nil {
			globalErr = err
			return err
		}
		_, err = connection.Exec(context.Background(), string(fileContent))

		if err != nil {
			globalErr = err
			return nil
		}

		return nil
	})
	return globalErr

}

// Db interface will be mocked once initial release is done

type PostgresDriver struct {
	conn *pgx.Conn
}

func (pd *PostgresDriver) Close(ctx context.Context) error {

	return pd.conn.Close(ctx)

}

// Connect to postgres instance
func CreatePostgresDriver(connectionURL string) (*PostgresDriver, error) {

	conn, err := pgx.Connect(context.Background(), connectionURL)
	if err != nil {
		return &PostgresDriver{}, err
	}

	// Attempt to migrate on creation of interface driver
	if err = tryMigrate(conn); err != nil {
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

// Get user from email[email is currently a unique identifier]
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

// directly update user credits
func (pd *PostgresDriver) UpdateUserCredits(ctx context.Context, email string, newcredits uint) error {

	_, err := pd.conn.Exec(ctx, `UPDATE usersrecord SET credits=$1 WHERE email=$2`, newcredits, email)
	return err

}

func (pd *PostgresDriver) NullifyUserCredits(ctx context.Context, email string) error {

	_, err := pd.conn.Exec(ctx, "UPDATE usersrecord SET credits=0 where email=$1", email)
	return err

}
