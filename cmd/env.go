package cmd

import (
	"github.com/joho/godotenv"
)

func loadEnvVars() error {

	return godotenv.Load(".env")

}
