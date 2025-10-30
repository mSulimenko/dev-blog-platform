package main

import (
	"context"
	"fmt"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/config"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/repository"
	authService "github.com/mSulimenko/dev-blog-platform/internal/auth/service"
	authGrpc "github.com/mSulimenko/dev-blog-platform/internal/auth/transport/grpc"
	httphandler "github.com/mSulimenko/dev-blog-platform/internal/auth/transport/http"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/transport/kafka"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/database"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/logger"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting auth")
	cfg := config.Load()

	log := logger.New(cfg.Env)
	defer log.Sync()

	// db
	dbpool, err := database.NewPool(context.Background(), cfg.DB.Dsn)
	if err != nil {
		log.Error("cannot connect to database: ", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	log.Info("Connected to database")

	err = database.RunMigrations(dbpool, cfg.DB.MigrationsDir)
	if err != nil {
		log.Error("migrations failed: ", err)
		os.Exit(1)
	}
	log.Infof("Migrations applied successfully")

	// repo
	usersRepo := repository.NewUsersRepository(dbpool)

	// Инициализируем Kafka Dispatcher
	kafkaDispatcher, err := kafka.NewKafkaDispatcher(cfg.Kafka.Brokers, log)
	if err != nil {
		log.Error("kafka initialization failed: ", err)
		os.Exit(1)
	}
	defer kafkaDispatcher.Close()

	// services
	userService := authService.NewUsersService(
		usersRepo,
		kafkaDispatcher,
		log,
		cfg.Auth.AccessSecret,
		cfg.Auth.AccessDuration,
	)

	// router
	handler := httphandler.NewHandler(userService, log)
	router := handler.InitRouter()

	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	// grpc server
	conn, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("tcp connection failed: %w", err)
	}
	gRPCServer := grpc.NewServer()
	authServ := authService.NewAuthService(usersRepo, log, cfg.Auth.AccessSecret)
	authGrpc.Register(gRPCServer, authServ)
	go func() {
		log.Infof("Starting grpc server on :%s", cfg.GRPC.Port)
		gRPCServer.Serve(conn)
	}()

	log.Infof("Starting server on %s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	if err = srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %w", err)
	}

}
