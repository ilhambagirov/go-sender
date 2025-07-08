package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"go-sender/internal/model/dto/message"
	"log"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 50 * time.Second,
}

type webhookSender struct {
	Url string
}

func NewWebhookSender(url string) ISender {
	return &webhookSender{Url: url}
}

func (w *webhookSender) Send(ctx context.Context, msg message.PayloadDto) (*message.ResponseDto, error) {
	log.Println("Sending message.")
	body, err := json.Marshal(msg)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.Url, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r message.ResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	log.Println("Message sent successfully.")
	return &r, nil
}
