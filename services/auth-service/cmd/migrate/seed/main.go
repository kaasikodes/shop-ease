package main

import (
	"log"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/db"
	"github.com/kaasikodes/shop-ease/shared/env"

	store "github.com/kaasikodes/shop-ease/services/auth-service/internal/store/sql-store"
	"github.com/kaasikodes/shop-ease/shared/logger"
)

func main() {
	addr := env.GetString("DB_ADDR", "root:root123$@tcp(localhost:3306)/shop_ease")
	logCfg := logger.LogConfig{
		LogFilePath: "../../logs/auth-service.log",
	}
	l := logger.New(logCfg)
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
