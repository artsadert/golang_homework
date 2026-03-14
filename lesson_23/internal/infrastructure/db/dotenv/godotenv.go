package dotenv

import (
	"github.com/joho/godotenv"
)

func LoadDotenv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
