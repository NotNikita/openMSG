package service

import (
	"context"
	"errors"

	"app/internal/models"
	"app/internal/repository"
)

type ConversationService interface {
	GetOrCreateConversation(ctx context.Context, userAID, userBID string) (*models.Conversation, error)
	ListConversations(ctx context.Context, userID string) ([]models.Conversation, error)
	GetConversation(ctx context.Context, id string) (*models.Conversation, error)
}

type conversationService struct {
	repo repository.ConversationRepository
}

func NewConversationService(repo repository.ConversationRepository) ConversationService {
	return &conversationService{repo: repo}
}

func (s *conversationService) GetOrCreateConversation(ctx context.Context, userAID, userBID string) (*models.Conversation, error) {
	if userAID == "" || userBID == "" {
		return nil, errors.New("user_a_id and user_b_id are required")
	}
	if userAID == userBID {
		return nil, errors.New("cannot create conversation with yourself")
	}
	// Sort to enforce the UNIQUE(user_a_id, user_b_id) constraint consistently
	if userAID > userBID {
		userAID, userBID = userBID, userAID
	}
	return s.repo.GetOrCreate(ctx, userAID, userBID)
}

func (s *conversationService) ListConversations(ctx context.Context, userID string) ([]models.Conversation, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	return s.repo.GetByUserID(ctx, userID)
}

func (s *conversationService) GetConversation(ctx context.Context, id string) (*models.Conversation, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}
