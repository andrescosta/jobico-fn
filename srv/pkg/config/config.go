package config

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/joho/godotenv"
)

var Environment string = DEVELOPMENT

const (
	DEVELOPMENT = "development"
	PRODUCTION  = "production"
	TEST        = "test"
)

var Environments = []string{DEVELOPMENT, PRODUCTION, TEST}

// Follows this convention: https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use
func LoadEnvVariables() error {
	Environment = os.Getenv("APP_ENV")
	if strings.TrimSpace(Environment) == "" {
		Environment = DEVELOPMENT
	} else {
		if !slices.Contains(Environments, Environment) {
			return fmt.Errorf("Invalid environment %s", Environment)
		}
	}

	godotenv.Load(".env." + Environment + ".local")
	if "test" != Environment {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + Environment)
	godotenv.Load()
	return nil
}
