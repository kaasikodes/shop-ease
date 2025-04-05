package config

import "github.com/kaasikodes/shop-ease/shared/env"

type Config struct {
	ApiAddr  string
	GrpcAddr string
	Env      string
	Db       DbConfig
	Mail     MailConfig
}

type MailConfig struct {
	FromEmail string
	Host      string
	Port      int
	Username  string
	Password  string
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
	Mail: MailConfig{
		Host: "sandbox.smtp.mailtrap.io",
		// Host:      "smtp.mailtrap.io",
		FromEmail: "hello@shop-ease.com",
		Port:      2525,
		Username:  "6c53d765680ca4",
		Password:  "83175273732073",
	},
}
