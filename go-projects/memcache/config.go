package memcache

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	MemoryThreshold   uint64
	DefaultExpiration time.Duration
	ServerPort        string
	ActiveTrace       bool
	ServerMode        ServeMode
	LogLevel          zerolog.Level
}

var Logger zerolog.Logger

type ServeMode string

const (
	TCP  ServeMode = "tcp"
	UDP  ServeMode = "udp"
	BOTH ServeMode = "both"
)

func NewConfig() *Config {
	var c Config

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	if c.MemoryThreshold, err = strconv.ParseUint(os.Getenv("LOCAL_MEMORY_THRESHOLD_IN_MBS"), 10, 64); err != nil {
		panic("LOCAL_MEMORY_THRESHOLD_IN_MBS is not set")
	}
	c.MemoryThreshold *= 1024 * 1024

	if exp, err := strconv.Atoi(os.Getenv("DEFAULT_EXPIRATION_IN_SECONDS")); err != nil {
		panic("DEFAULT_EXPIRATION_IN_SECONDS is not set")

	} else {
		c.DefaultExpiration = time.Duration(exp) * time.Second
	}

	if c.ServerPort = os.Getenv("SERVER_PORT"); c.ServerPort == "" {
		panic("SERVER_PORT is not set")

	}

	if c.ActiveTrace, err = strconv.ParseBool(os.Getenv("ACTIVE_TRACE")); err != nil {
		panic("ACTIVE_TRACE is not set")

	}

	switch os.Getenv("SERVER_MODE") {
	case "TCP":
		c.ServerMode = TCP
	case "UDP":
		c.ServerMode = UDP
	case "BOTH":
		c.ServerMode = BOTH
	default:
		panic("SERVER_MODE is not set properly. Valid values are TCP, UDP, BOTH")
	}

	if c.LogLevel, err = zerolog.ParseLevel(os.Getenv("LOG_LEVEL")); err != nil {
		panic("LOG_LEVEL is not set properly. Valid values are debug, info, warn, error, fatal, panic")
	}

	setupLogger(c.LogLevel)
	return &c
}

func setupLogger(level zerolog.Level) {
	runLogFile, err := os.OpenFile(
		"log.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic("Error opening the log file")
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = "02-01-2006 à 15:04:05"
	output := zerolog.MultiLevelWriter(os.Stdout, runLogFile)
	Logger = zerolog.New(output).With().Timestamp().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
