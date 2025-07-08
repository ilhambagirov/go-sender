package util

import (
	"fmt"
	"go-sender/internal/config"
	"go-sender/internal/service/sender"
	"math/rand"
)

func RandomString(size int) string {
	symbols := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, size)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

func RandomPhone(size int) string {
	symbols := "0123456789"
	b := make([]byte, size)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

func GetSender() (sender.ISender, error) {
	senderType := config.Config.SenderType
	switch senderType {
	case "webhook":
		return sender.NewWebhookSender(config.Config.WebhookUrl), nil
	default:
		return nil, fmt.Errorf("unknown sender type %q", senderType)
	}
}
