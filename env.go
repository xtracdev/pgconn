package pgconn

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type EnvConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

const (
	DBUser     = "DB_USER"
	DBPassword = "DB_PASSWORD"
	DBHost     = "DB_HOST"
	DBPort     = "DB_PORT"
	DBName     = "DB_NAME"
)

func (ec *EnvConfig) MaskedConnectString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		ec.DBUser, "XXX", ec.DBName, ec.DBHost, ec.DBPort)
}

func (ec *EnvConfig) ConnectString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		ec.DBUser, ec.DBPassword, ec.DBName, ec.DBHost, ec.DBPort)
}

func NewEnvConfig() (*EnvConfig, error) {
	var configErrors []string

	user := os.Getenv(DBUser)
	if user == "" {
		configErrors = append(configErrors, "Configuration missing DB_USER env variable")
	}

	password := os.Getenv(DBPassword)
	if password == "" {
		configErrors = append(configErrors, "Configuration missing DB_PASSWORD env variable")
	}

	dbhost := os.Getenv(DBHost)
	if dbhost == "" {
		configErrors = append(configErrors, "Configuration missing DB_HOST env variable")
	}

	dbPort := os.Getenv(DBPort)
	if dbPort == "" {
		configErrors = append(configErrors, "Configuration missing DB_PORT env variable")
	}

	dbSName := os.Getenv(DBName)
	if dbSName == "" {
		configErrors = append(configErrors, "Configuration missing DB_NAME env variable")
	}

	if len(configErrors) != 0 {
		return nil, errors.New(strings.Join(configErrors, "\n"))
	}

	return &EnvConfig{
		DBUser:     user,
		DBPassword: password,
		DBHost:     dbhost,
		DBPort:     dbPort,
		DBName:     dbSName,
	}, nil

}
