package main

import (
	"github.com/clearcodecn/carim/server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	s := server.NewHttpServer(server.Config{
		Driver: "sqlite3",
		DSN:    "./car.db",
		Key:    "car",

		Host: "0.0.0.0",
		Port: "9527",
	})

	s.Start()
}
