package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
	"net/url"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

func NewConnPool(dbConfig DBConfig) (*pgxpool.Pool, error) {
	databaseUrl := formUrl("postgres", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	dbPool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, err
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return dbPool, err
}

func formUrl(scheme, username, password, host, port, path string) string {
	var u = url.URL{
		Scheme: scheme,
		User:   url.UserPassword(username, password),
		Host:   net.JoinHostPort(host, port),
		Path:   path,
	}
	return u.String()
}
