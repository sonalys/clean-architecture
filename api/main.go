package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spoonrocker/cart-go-sonalys/pkg/config"
	"github.com/spoonrocker/cart-go-sonalys/pkg/http"
	"github.com/spoonrocker/cart-go-sonalys/pkg/persistence"
	"github.com/spoonrocker/cart-go-sonalys/usecase"
)

type (
	Services struct {
		http http.HTTP
		db   persistence.Persistence
	}

	APIConfig struct {
		HTTP struct {
			Address   string `yaml:"host_address"`
			RateLimit int    `yaml:"rate_limit"`
			Metrics   struct {
				Route   string `yaml:"route"`
				Address string `yaml:"host_address"`
			} `yaml:"metrics"`
		} `yaml:"http"`

		Persistence struct {
			User         string `yaml:"user"`
			Password     string `yaml:"password"`
			Host         string `yaml:"host"`
			Port         string `yaml:"port"`
			DatabaseName string `yaml:"database_name"`
			SSL          bool   `yaml:"ssl"`
		} `yaml:"database"`
	}
)

func ReadConfig() (APIConfig, error) {
	c := new(APIConfig)
	configPath := config.GetEnvVar("CONFIG_PATH", "./api/config.yaml")

	err := config.ReadYAML(configPath, c)
	if err != nil {
		return APIConfig{}, err
	}

	return *c, nil
}

func main() {
	logrus.Info("initializing API")

	c, err := ReadConfig()
	if err != nil {
		logrus.Error("failed to read config")
		return
	}

	postgresConn := persistence.BuildDatabaseConnString(
		c.Persistence.User,
		c.Persistence.Password,
		c.Persistence.Host,
		c.Persistence.Port,
		c.Persistence.DatabaseName,
		c.Persistence.SSL)

	db, err := persistence.NewPersistence(postgresConn)
	if err != nil {
		logrus.Error("failed to initialize the database service")
		return
	}

	s := Services{
		http: http.CreateHTTP(c.HTTP.RateLimit),
		db:   db,
	}

	usecase.CreateCartUsecase(s.http, s.db)

	s.http.ListenMetrics(c.HTTP.Metrics.Route, c.HTTP.Metrics.Address)
	logrus.Fatal(s.http.Listen(c.HTTP.Address))
}
