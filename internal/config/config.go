package config

import (
	"errors"
	"fmt"
	"net"

	"github.com/UdinSemen/golang-test-task/pkg/connectors/postgres"
	"github.com/go-playground/validator/v10"
)

// App
// struct of application config
type App struct {
	Connectors struct {
		HTTPServer  HTTPServer `env-prefix:"HTTP_SERVER_" yaml:"httpServer"`
		MainStorage struct {
			Postgres postgres.DatabaseConfig `yaml:"postgres" env-prefix:"POSTGRES_"`
		} `yaml:"mainStorage" env-prefix:"MAIN_STORAGE_"`
	} `yaml:"connectors" env-prefix:"CONNECTORS_"`
}

func (a *App) Validate() error {
	err := validator.New().Struct(a)
	var errJoin error
	if err != nil {
		var validationError validator.ValidationErrors
		if errors.As(err, &validationError) {
			for _, field := range validationError {
				errJoin = errors.Join(err, fmt.Errorf(
					"field = %s, namespace = %s, tag = %s",
					field.Field(),
					field.Namespace(),
					field.Tag(),
				))
			}
		}
	}

	return errJoin
}

type HTTPServer struct {
	Host string `env:"HOST" yaml:"host" validate:"required"`
	Port string `env:"PORT" yaml:"port" validate:"required"`
}

func (c HTTPServer) GetAddress() string {
	return net.JoinHostPort(c.Host, c.Port)
}
