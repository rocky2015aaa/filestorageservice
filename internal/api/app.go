package api

import (
	"log"
	"net/http"
	"os"

	"github.com/rocky2015aaa/filestorageservice/internal/api/handlers"
	"github.com/rocky2015aaa/filestorageservice/internal/config"
	"github.com/rocky2015aaa/filestorageservice/internal/pkg/database"
)

func NewApp(sig chan os.Signal) *http.Server {
	db, err := database.NewPostgres(os.Getenv(config.EnvSvrPostgresUri))
	if err != nil {
		log.Fatalln(err)
	}
	return &http.Server{
		Addr:    ":" + os.Getenv(config.EnvSvrPort),
		Handler: NewRouter(handlers.NewHandler(db)),
	}
}
