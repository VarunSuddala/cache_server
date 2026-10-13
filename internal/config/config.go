package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
	HOST string
}

func MustLoad() Config {
	godotenv.Load()
	port:= os.Getenv("PORT")
	if port == ""{
		panic("PORT is required")
	}
	host := os.Getenv("HOST")
	if host == ""{
		panic("HOST is required")
	}

	return Config{
		PORT: port,
		HOST: host,
	}

}