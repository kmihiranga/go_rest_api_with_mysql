package main

import (
	"context"
	"database/sql"
	"fmt"
	"go_rest_api_with_mysql/api/handler"
	"go_rest_api_with_mysql/config"
	"go_rest_api_with_mysql/infrastructure/repository"
	logger "go_rest_api_with_mysql/pkg/log"
	"go_rest_api_with_mysql/usecase/user"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gorillaCtx "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger = logger.GetLogger().Sugar()

func main() {

	cfg := config.Parse()
	// ctx := context.Background()
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.Database)
	conn, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	// verify the connection
	if err := conn.Ping(); err != nil {
		log.Fatalf("unable to reach databases. %v", err.Error())
	}

	userRepository := repository.NewUserMySQL(conn)

	userService := user.NewService(cfg.AppConfig, userRepository)

	// router configs
	r := mux.NewRouter()
	http.Handle("/", r)

	// heartbeat
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.RegisterUserHandler(r, userService)
	handler.RegisterNotFoundHandler(r)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.AppConfig.Host, cfg.AppConfig.Port),
		ReadTimeout:  time.Duration(cfg.AppConfig.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.AppConfig.WriteTimeout) * time.Millisecond,
		Handler:      gorillaCtx.ClearHandler(http.DefaultServeMux),
	}

	// using goroutine to graceful shutdown the server
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	// deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.AppConfig.ShutdownWaitTimeout)*time.Millisecond)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Warnf("Server Shutdown Failed. %+v", err)
	}

	log.Info("Server shutdown completed.")

	os.Exit(0)
}
