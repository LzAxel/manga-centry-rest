package main

import (
	"mangacentry"
	"mangacentry/pkg/handler"
	"mangacentry/pkg/repository"
	"mangacentry/pkg/service"

	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"

	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.Infoln("Starting")
	if err := initConfig(); err != nil {
		logrus.Fatalf("can't load config file: %s", err.Error())
	}
	if level := viper.GetBool("debug"); level {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("debug mode enabled")
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("can't connect to database: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(*services)

	server := new(mangacentry.Server)
	if err := server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("server not started: %s", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
