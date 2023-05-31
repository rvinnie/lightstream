package postgres

import (
	"context"
	"net"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

func NewConnPool(dbConfig DBConfig) (*pgxpool.Pool, error) {
	databaseUri := formUri(dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	dbPool, err := pgxpool.New(context.Background(), databaseUri)
	if err != nil {
		return nil, err
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return dbPool, err
}

func formUri(username, password, host, port, dbname string) string {
	var u = url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Host:   net.JoinHostPort(host, port),
		Path:   dbname,
	}
	return u.String()
}
