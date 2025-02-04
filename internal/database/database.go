package database

import "context"

type DBConn interface {
	Ping(context.Context) error
	Close() error
	Query(query string) QueryResult
	QueryRows(query string) QueryResult
	Select(query string)
}

type QueryResult interface {
	Scan(scanto any) error
}
