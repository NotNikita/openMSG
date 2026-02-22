package repository

import (
	"context"

	"app/internal/models"
	"app/internal/repository/sqlcgen"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate minimock -i ConversationRepository -o ../mocks/conversation_repo_mock.go -n ConversationRepositoryMock -p mocks

type ConversationRepository interface {
	GetOrCreate(ctx context.Context, userAID, userBID string) (*models.Conversation, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Conversation, error)
	GetByID(ctx context.Context, id string) (*models.Conversation, error)
}

type postgresConversationRepository struct {
	q *sqlcgen.Queries
}

func NewConversationRepository(pool *pgxpool.Pool) ConversationRepository {
	return &postgresConversationRepository{q: sqlcgen.New(pool)}
}

func (r *postgresConversationRepository) GetOrCreate(ctx context.Context, userAID, userBID string) (*models.Conversation, error) {
	row, err := r.q.GetOrCreateConversation(ctx, &sqlcgen.GetOrCreateConversationParams{
		UserAID: userAID,
		UserBID: userBID,
	})
	if err != nil {
		return nil, err
	}
	return toConversation(row), nil
}

func (r *postgresConversationRepository) GetByUserID(ctx context.Context, userID string) ([]models.Conversation, error) {
	rows, err := r.q.GetConversationsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	convs := make([]models.Conversation, len(rows))
	for i, row := range rows {
		convs[i] = *toConversation(row)
	}
	return convs, nil
}

func (r *postgresConversationRepository) GetByID(ctx context.Context, id string) (*models.Conversation, error) {
	row, err := r.q.GetConversationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}
	return toConversation(row), nil
}

func toConversation(r *sqlcgen.Conversation) *models.Conversation {
	return &models.Conversation{
		ID:        r.ID,
		UserAID:   r.UserAID,
		UserBID:   r.UserBID,
		CreatedAt: r.CreatedAt.Time,
	}
}
