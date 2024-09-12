package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	Date    = ""
	Version = ""
	Build   = ""
)

const (
	EnvSvrPort            = "SVR_PORT"
	EnvSvrLogLevel        = "SVR_LOG_LEVEL"
	EnvSvrGinMode         = "SVR_GIN_MODE"
	EnvSvrPostgresUri     = "SVR_POSTGRES_URI"
	EnvSvrFileSliceNumber = "SVR_FILE_SLICE_NUMBER"

	EnvFile = ".env"
)

func init() {
	fmt.Printf("Build Date: %s\nBuild Version: %s\nBuild: %s\n\n", Date, Version, Build)
	err := godotenv.Load(EnvFile)
	if err != nil {
		log.Fatalf("Error loading %s file: %v", EnvFile, err)
	}
	logLevel, err := log.ParseLevel(os.Getenv(EnvSvrLogLevel))
	if err != nil {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)
	log.SetFormatter(&log.JSONFormatter{})
}
