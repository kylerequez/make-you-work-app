package utils

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	log.Println("::: Loading Env Variables...")
	if err := godotenv.Load(); err != nil {
		return err
	}
	log.Println("::: Successfully Loaded Env Variables")
	return nil
}
