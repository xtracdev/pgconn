package pgconn

import (
	"errors"
	"fmt"
	"github.com/xtracdev/envinject"
	"strings"
)

type envConfig struct {
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

func (ec *envConfig) maskedConnectString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		ec.DBUser, "XXX", ec.DBName, ec.DBHost, ec.DBPort)
}

func (ec *envConfig) connectString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		ec.DBUser, ec.DBPassword, ec.DBName, ec.DBHost, ec.DBPort)
}

func newEnvConfig(env *envinject.InjectedEnv) (*envConfig, error) {
	var configErrors []string

	user := env.Getenv(DBUser)
	if user == "" {
		configErrors = append(configErrors, "Configuration missing DB_USER env variable")
	}

	password := env.Getenv(DBPassword)
	if password == "" {
		configErrors = append(configErrors, "Configuration missing DB_PASSWORD env variable")
	}

	dbhost := env.Getenv(DBHost)
	if dbhost == "" {
		configErrors = append(configErrors, "Configuration missing DB_HOST env variable")
	}

	dbPort := env.Getenv(DBPort)
	if dbPort == "" {
		configErrors = append(configErrors, "Configuration missing DB_PORT env variable")
	}

	dbSName := env.Getenv(DBName)
	if dbSName == "" {
		configErrors = append(configErrors, "Configuration missing DB_NAME env variable")
	}

	if len(configErrors) != 0 {
		return nil, errors.New(strings.Join(configErrors, "\n"))
	}

	return &envConfig{
		DBUser:     user,
		DBPassword: password,
		DBHost:     dbhost,
		DBPort:     dbPort,
		DBName:     dbSName,
	}, nil

}

func ConnectStringFromInjectedEnv(env *envinject.InjectedEnv) (string, error) {
	if env == nil {
		return "", errors.New("Nil InjectedEnv")
	}

	config, err := newEnvConfig(env)
	if err != nil {
		return "", err
	}

	return config.connectString(), nil
}

func MaskedConnectStringFromInjectedEnv(env *envinject.InjectedEnv) (string, error) {
	if env == nil {
		return "", errors.New("Nil InjectedEnv")
	}

	config, err := newEnvConfig(env)
	if err != nil {
		return "", err
	}

	return config.maskedConnectString(), nil
}
