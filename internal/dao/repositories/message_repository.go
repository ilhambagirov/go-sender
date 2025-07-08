package repositories

import (
	"go-sender/internal/dao/domain"
	"go-sender/internal/model/dto/message"
	"go-sender/internal/util"
	"gorm.io/gorm"
)

type IMessageRepository interface {
	CreateMessage(request *domain.Message) error
	GetMessages(isSent bool, paging *util.Paging) ([]domain.Message, error)
	UpdateMessage(msg message.Dto) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) IMessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) CreateMessage(request *domain.Message) error {
	result := r.db.Create(request)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *messageRepository) GetMessages(isSent bool, paging *util.Paging) ([]domain.Message, error) {
	var messages []domain.Message
	result := r.db.Scopes(util.NewPaging(paging).PaginatedResult).
		Where("is_sent = ?", isSent).
		Order("created_at ASC").
		Find(&messages)

	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

func (r *messageRepository) UpdateMessage(msg message.Dto) error {
	updated := domain.Message{
		IsSent: msg.IsSent,
	}
	result := r.db.
		Where("id = ?", msg.Id).
		Updates(updated)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
