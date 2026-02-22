package service

import (
	"context"
	"errors"

	"app/internal/models"
	"app/internal/repository"
)

type MessageService interface {
	SendMessage(ctx context.Context, conversationID, senderID, ciphertext, nonce string) (*models.Message, error)
	ListMessages(ctx context.Context, conversationID string) ([]models.Message, error)
	ListPublicMessages(ctx context.Context) ([]models.PublicMessage, error)
}

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) MessageService {
	return &messageService{repo: repo}
}

func (s *messageService) SendMessage(ctx context.Context, conversationID, senderID, ciphertext, nonce string) (*models.Message, error) {
	if conversationID == "" {
		return nil, errors.New("conversation_id is required")
	}
	if senderID == "" {
		return nil, errors.New("sender_id is required")
	}
	if ciphertext == "" {
		return nil, errors.New("ciphertext is required")
	}
	if nonce == "" {
		return nil, errors.New("nonce is required")
	}
	return s.repo.Create(ctx, conversationID, senderID, ciphertext, nonce)
}

func (s *messageService) ListMessages(ctx context.Context, conversationID string) ([]models.Message, error) {
	if conversationID == "" {
		return nil, errors.New("conversationID is required")
	}
	return s.repo.GetByConversationID(ctx, conversationID)
}

func (s *messageService) ListPublicMessages(ctx context.Context) ([]models.PublicMessage, error) {
	return s.repo.GetPublic(ctx)
}
