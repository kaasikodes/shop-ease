package config

import "github.com/kaasikodes/shop-ease/shared/env"

type Config struct {
	ApiAddr  string
	GrpcAddr string
	Env      string
	Db       DbConfig
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

var ServiceConfig = Config{

	ApiAddr:  env.GetString("API_ADDR", ":3020"),
	GrpcAddr: env.GetString("GRPC_ADDR", ":5050"),
	Db: DbConfig{
		Addr:         env.GetString("DB_ADDR", "mysql://root:root123$@localhost"),
		MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	},
}
