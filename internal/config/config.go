package config

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"go-sender/internal/dao/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type args struct {
	Port       int    `arg:"env:PORT"`
	DbHost     string `arg:"env:DB_HOST"`
	DbPort     string `arg:"env:DB_PORT"`
	DbUser     string `arg:"env:POSTGRES_USER"`
	DbName     string `arg:"env:POSTGRES_DB"`
	DbPassword string `arg:"env:POSTGRES_PASSWORD"`
	WebhookUrl string `arg:"env:WEBHOOK_URL"`
	RedisHost  string `arg:"env:REDIS_HOST"`
	RedisPort  string `arg:"env:REDIS_PORT"`
	RedisPass  string `arg:"env:REDIS_PASSWORD"`
	RedisDb    string `arg:"env:REDIS_DB"`
	SenderType string `arg:"env:SENDER_TYPE"`
}

var Config args

func Load() {
	_ = arg.MustParse(&Config)
}

func InitDb() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Config.DbHost, Config.DbUser, Config.DbPassword, Config.DbName, Config.DbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = db.AutoMigrate(&domain.Message{})
	if err != nil {
		panic(err)
	}

	return db
}
