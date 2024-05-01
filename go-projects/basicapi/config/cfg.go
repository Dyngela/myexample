package config

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"

	"time"
)

type Mode string

const (
	Dev  Mode = "dev"
	Prod Mode = "prod"
)

type Configuration struct {
	ServerPort string
	ServerMode string

	PGHost     string
	PGPort     string
	PGUsername string
	PGPassword string
	PGDbname   string

	ApiKey string
}

var Logger zerolog.Logger
var DB *gorm.DB
var config *Configuration

func InitConfiguration() {
	InitLogger()

	if err := godotenv.Load(".env"); err != nil {
		Logger.Error().Msgf("Error loading .env file: %s Falling back into default mode value", err.Error())
	}

	config = &Configuration{
		ServerPort: os.Getenv("SERVER_PORT"),
		ServerMode: os.Getenv("SERVER_MODE"),
		PGHost:     os.Getenv("PG_HOST"),
		PGPort:     os.Getenv("PG_PORT"),
		PGUsername: os.Getenv("PG_USERNAME"),
		PGPassword: os.Getenv("PG_PASSWORD"),
		PGDbname:   os.Getenv("PG_DBNAME"),
		ApiKey:     os.Getenv("API_KEY"),
	}

	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}
	if config.ServerMode == "" {
		config.ServerMode = "dev"
	}
	if config.PGHost == "" {
		config.PGHost = "localhost"
	}
	if config.PGPort == "" {
		config.PGPort = "5432"
	}
	if config.PGUsername == "" {
		config.PGUsername = "postgres"
	}
	if config.PGPassword == "" {
		config.PGPassword = "postgres"
	}
	if config.PGDbname == "" {
		config.PGDbname = "basicapi"
	}
	if config.ApiKey == "" {
		config.ApiKey = "basicapi"
	}

	ConnectToPostgres()
}

func GetConfig() *Configuration {
	return config
}

func ConnectToPostgres() {
	var err error
	var db *gorm.DB
	var conn *sql.DB

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.PGHost, config.PGUsername, config.PGPassword, config.PGDbname, config.PGPort)
	if db, err = gorm.Open(postgres.Open(dsn),
		&gorm.Config{
			TranslateError: true,
			NowFunc: func() time.Time {
				//loc, err := time.LoadLocation("Europe/Paris")
				//if err != nil {
				//	panic(err)
				//}
				// No idea why but it's not working without adding 1 hour
				//return time.Now().In(loc).Add(time.Hour * 1)
				//return time.Now().Add(time.Hour * 1)
				return time.Now()
			},
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:         "",
				SingularTable:       true,
				NameReplacer:        nil,
				NoLowerCase:         false,
				IdentifierMaxLength: 0,
			}}); err != nil {
	}
	if conn, err = db.DB(); err != nil {
		panic(err)
	}
	conn.SetMaxIdleConns(10)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)
	DB = db
}

func InitLogger() {
	runLogFile, _ := os.OpenFile(
		"basicapi.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	zerolog.TimeFieldFormat = "02-01-2006 à 15:04:05"
	output := zerolog.MultiLevelWriter(os.Stdout, runLogFile)
	Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
}
