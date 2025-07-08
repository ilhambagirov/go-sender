package sender

import (
	"context"
	"go-sender/internal/model/dto/message"
)

type ISender interface {
	Send(ctx context.Context, payload message.PayloadDto) (*message.ResponseDto, error)
}
