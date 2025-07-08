package service

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go-sender/internal/dao/domain"
	"go-sender/internal/dao/repositories"
	"go-sender/internal/model/dto/message"
	"go-sender/internal/service/sender"
	"go-sender/internal/util"
	"log"
	"math/rand"
	"time"
)

type IMessageService interface {
	CreateMessage(request message.CreateMessageRequest) error
	GetMessages(isSent bool, paging *util.Paging) ([]message.Dto, error)
	Initialize()
	Start() error
	Stop()
}

type messageService struct {
	mr     repositories.IMessageRepository
	sender sender.ISender
	redis  *redis.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMessageService(mr repositories.IMessageRepository, sender sender.ISender,
	ctx context.Context, redis *redis.Client) IMessageService {
	return &messageService{
		mr:     mr,
		sender: sender,
		redis:  redis,
		ctx:    ctx,
		cancel: nil,
	}
}

func (s *messageService) CreateMessage(request message.CreateMessageRequest) error {
	msg := &domain.Message{
		Phone:   request.Content,
		Content: request.Content,
		IsSent:  false,
	}

	err := s.mr.CreateMessage(msg)

	if err != nil {
		return err
	}

	return nil
}

func (s *messageService) GetMessages(isSent bool, paging *util.Paging) ([]message.Dto, error) {
	messages, err := s.mr.GetMessages(isSent, paging)
	if err != nil {
		return nil, err
	}

	msgDtos := make([]message.Dto, 0)

	for _, msg := range messages {
		msgDtos = append(msgDtos, message.Dto{
			Id:          msg.ID,
			Content:     msg.Content,
			Phone:       msg.Phone,
			IsSent:      msg.IsSent,
			CreatedDate: msg.CreatedAt,
		})
	}

	return msgDtos, nil
}

func (s *messageService) Initialize() {
	if s.cancel != nil {
		return
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
}

func (s *messageService) Start() error {
	if s.ctx == nil {
		return errors.New("service not initialized")
	}

	go func() {
		if err := s.process(); err != nil {
			log.Printf("process exited: %v", err)
		}
	}()

	go func() {
		if err := s.insertNewMessage(); err != nil {
			log.Printf("insertNewMessage exited: %v", err)
		}
	}()
	return nil
}

func (s *messageService) Stop() {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}

// process starts sending messages
func (s *messageService) process() error {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case <-ticker.C:
			if err := s.processBatch(s.ctx); err != nil {
				return err
			}
		}
	}
}

// insertNewMessage writes new messages every 2 seconds for newly added fresh records
func (s *messageService) insertNewMessage() error {
	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case <-ticker.C:
			msg := &domain.Message{
				Content: util.RandomString(rand.Intn(100)),
				Phone:   util.RandomPhone(12),
				IsSent:  false,
			}
			if err := s.mr.CreateMessage(msg); err != nil {
				return err
			}
		}
	}
}

// processBatch prepares messages, calls send and updates message.
func (s *messageService) processBatch(ctx context.Context) error {
	messages, err := s.GetMessages(false, &util.Paging{Limit: 2, Page: 1})
	if err != nil {
		return err
	}

	for _, msg := range messages {
		readyMsg := message.PayloadDto{
			Content: msg.Content,
			To:      msg.Phone,
		}
		res, err := s.sender.Send(ctx, readyMsg)
		if err != nil {
			log.Printf("failed to send message %v: %v", msg.Id, err)
			continue
		}
		sentAt := time.Now()

		msg.IsSent = true
		if err := s.mr.UpdateMessage(msg); err != nil {
			log.Printf("failed to update message %v: %v", msg.Id, err)
		}

		// cache sent message
		s.redis.Set(s.ctx, res.MessageId, sentAt, 24*time.Hour)
	}
	return nil
}
