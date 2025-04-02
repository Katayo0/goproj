package main

import (
	"fmt"
	"go_sql_test/internal/config"
	"go_sql_test/internal/http-server/handlers/url/save"
	"go_sql_test/internal/http-server/middleware/logger"
	"go_sql_test/internal/lib/logger/sl"
	"go_sql_test/internal/storage/psql"
	"log/slog"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sumit-tembe/gin-requestid"

	_ "github.com/lib/pq" // init pq driver
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)


func main(){
	cfg := config.MustLoad()

	fmt.Println(nil)

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := psql.New(cfg.StorageUrl)
	if err != nil{
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	//new engine 
	engine := gin.New()

	//middleware
	engine.Use(requestid.RequestID(nil))
	engine.Use(logger.New(log))
	engine.Use(gin.Recovery())

	engine.POST("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      engine,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger{

	var log *slog.Logger

	switch env {

	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))	

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))	

	}
	
	return log
}