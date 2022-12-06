package config

import (
	"os"
	"strconv"
)

var (
	RPCPort = ":8787"

	JwtSecretKey string

	DbHost    = "localhost"
	DbPort    = "3030"
	DbTz      = "UTC"
	DbSslMode = "disable"

	DbUser     string
	DbPassword string
	DbName     string

	Debug = false
)

func init() {
	if x, err := strconv.ParseBool(os.Getenv("DEBUG")); err == nil {
		Debug = x
	}

	if envRpcPort := os.Getenv("RPC_PORT"); envRpcPort != "" {
		RPCPort = envRpcPort
	}

	if jwtKey := os.Getenv("JWT_SECRET_KEY"); jwtKey != "" {
		JwtSecretKey = jwtKey
	}

	if x := os.Getenv("DB_HOST"); x != "" {
		DbHost = x
	}

	if x := os.Getenv("DB_PORT"); x != "" {
		DbPort = x
	}

	if x := os.Getenv("DB_TZ"); x != "" {
		DbTz = x
	}

	if os.Getenv("DB_SSL_MODE") == "enable" {
		DbSslMode = "enable"
	}

	DbUser = os.Getenv("DB_USER")
	DbName = os.Getenv("DB_NAME")
	DbPassword = os.Getenv("DB_PASSWORD")
}
