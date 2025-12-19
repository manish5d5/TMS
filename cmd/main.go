package main

import (
	"context"
	"fmt"

	"log"
	"os"

	"os/signal"
	"syscall"

	"TMS/config"
	"TMS/https"
	"TMS/https/handler"
	"TMS/repos"
	"TMS/repos/postgres"
	service "TMS/services"
	"TMS/utils/jwt"

	"github.com/redis/go-redis/v9"
)

func InitializeServer(config config.Config, ctx context.Context) (*https.Server, error) {

	//databse

	pool, err := postgres.Connect(ctx, config)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect to redis: %v", err)
		return nil, err
	}
	secret := jwt.LoadJWT()

	//db
	userRepo := repos.NewUserRepo(pool)

	//services
	userService := service.NewUserService(userRepo)

	authService := service.NewAuthService(userRepo, &jwt.Jwt{
		SecretKey:     secret,
		SecretKeyByte: []byte(secret),
	})

	authHandler := handler.NewAuthHandler(authService, userService)
	userHandler := handler.NewUserHandler(userService)

	fmt.Println(pool)
	appserver := https.NewServer(config, authHandler, userHandler)

	return appserver, nil

}

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.DefaultConfig
	log.Println("initilaizng server")
	appserver, err := InitializeServer(cfg, ctx)
	if err != nil {
		log.Printf("failed ")
		return
	}
	go func() {
		<-ctx.Done()
		log.Println("Graceful shutdown initiated...")
	}()

	log.Printf("initialized on port %s", cfg.Listen)
	if err = appserver.Listen(ctx, cfg.Listen); err != nil {
		log.Println("error while connecting to server")
	}
}
