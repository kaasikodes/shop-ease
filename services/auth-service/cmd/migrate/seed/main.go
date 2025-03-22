package main

import (
	"log"

	"github.com/kaasikodes/shop-ease/cmd/logger"
	"github.com/kaasikodes/shop-ease/internal/db"
	"github.com/kaasikodes/shop-ease/internal/env"
	store "github.com/kaasikodes/shop-ease/internal/store/sql-store"
)

func main() {
	addr := env.GetString("DB_ADDR", "root:root123$@tcp(localhost:3306)/shop_ease")

	l := logger.New("../../app.log")
	l.Info(addr, "DB_ADDR ...")

	conn, err := db.New(addr, 10, 10, "5m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	s := store.NewSQLStorage(conn)
	err = db.Seed(s, conn)
	if err != nil {
		panic(err)
	}

}
