package main

import (
	"log"

	"github.com/kaasikodes/shop-ease/cmd/logger"
	"github.com/kaasikodes/shop-ease/internal/db"
	"github.com/kaasikodes/shop-ease/internal/env"
	"github.com/kaasikodes/shop-ease/internal/store"
)

// TODO: Refactor seed & unseed to avoid logic repetition
func main() {
	addr := env.GetString("DB_ADDR", "root:root123$@tcp(localhost:3306)/shop_ease")
	
	l := logger.New("../../app.log")
	l.Info(addr, "DB_ADDR ...", )

	conn, err := db.New(addr, 10, 10, "5m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	s := store.NewStorage(conn)
	err = db.UnSeed(s,conn)
	if err != nil {
		panic(err)
	}

}