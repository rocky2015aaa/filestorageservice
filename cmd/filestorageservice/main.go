package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/rocky2015aaa/filestorageservice/docs"
	"github.com/rocky2015aaa/filestorageservice/internal/api"
	log "github.com/sirupsen/logrus"
)

// @title   File Storage Service
// @version test

// @contact.name  Donggeon Lee
// @contact.email rocky2010aaa@gmail.com

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost
// @BasePath /
func main() {
	signalReceived := make(chan os.Signal, 1)
	svr := api.NewApp(signalReceived)
	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	signalReceived <- <-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdownRelease()

	if err := svr.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Server has shut down.")
}
