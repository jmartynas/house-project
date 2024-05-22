package controller

import (
	"context"
	"fmt"
	"main/database"

	"github.com/jackc/pgx/v5"
)

type Server struct {
	Host   string
	Port   string
	DBConn map[string]DBConn
}

type DBConn struct {
	user string
	host string
	port string
}

func NewServer(host, port string) *Server {
	return &Server{
		Host:   host,
		Port:   port,
		DBConn: make(map[string]DBConn),
	}
}

func NewDBConn(user, host, port string) DBConn {
	return DBConn{
		user: user,
		host: host,
		port: port,
	}
}

func (dbconn *DBConn) DBConn(ctx context.Context) (*database.Queries, error) {
	connurl := fmt.Sprintf("postgres://%v:%v@%v:%v/house", dbconn.user, "", dbconn.host, dbconn.port)
	conn, err := pgx.Connect(ctx, connurl)
	if err != nil {
		return &database.Queries{}, err
	}
	return database.New(conn), nil
}
