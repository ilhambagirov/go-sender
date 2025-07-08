package main

import (
	"context"
	"go-sender/internal/config"
	"go-sender/internal/dao/repositories"
	"go-sender/internal/handler"
	"go-sender/internal/service"
	"go-sender/internal/util"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config.Load()
	db := config.InitDb()
	redis, err := config.NewRedisClient()
	sender, err := util.GetSender()
	if err != nil {
		log.Fatal(err)
	}

	msgRepo := repositories.NewMessageRepository(db)
	msgService := service.NewMessageService(msgRepo, sender, ctx, redis)
	msgController := handler.NewMsgController(msgService, ctx)

	// automatic start sender
	if err := msgController.StartService(); err != nil {
		log.Fatalf("failed to start message service: %v", err)
	}

	handler.StartHttpServing(msgController)
}
