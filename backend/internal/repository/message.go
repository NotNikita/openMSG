package repository

import (
	"context"
	"time"

	"app/internal/models"
	"app/internal/repository/sqlcgen"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate minimock -i MessageRepository -o ../mocks/message_repo_mock.go -n MessageRepositoryMock -p mocks

type MessageRepository interface {
	Create(ctx context.Context, conversationID, senderID, ciphertext, nonce string) (*models.Message, error)
	GetByConversationID(ctx context.Context, conversationID string) ([]models.Message, error)
	GetPublic(ctx context.Context, limit int, before *time.Time) ([]models.PublicMessage, error)
}

type postgresMessageRepository struct {
	q *sqlcgen.Queries
}

func NewMessageRepository(pool *pgxpool.Pool) MessageRepository {
	return &postgresMessageRepository{q: sqlcgen.New(pool)}
}

func (r *postgresMessageRepository) Create(ctx context.Context, conversationID, senderID, ciphertext, nonce string) (*models.Message, error) {
	row, err := r.q.CreateMessage(ctx, &sqlcgen.CreateMessageParams{
		ConversationID: conversationID,
		SenderID:       senderID,
		Ciphertext:     ciphertext,
		Nonce:          nonce,
	})
	if err != nil {
		return nil, err
	}
	return toMessage(row), nil
}

func (r *postgresMessageRepository) GetByConversationID(ctx context.Context, conversationID string) ([]models.Message, error) {
	rows, err := r.q.GetMessagesByConversationID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	msgs := make([]models.Message, len(rows))
	for i, row := range rows {
		msgs[i] = *toMessage(row)
	}
	return msgs, nil
}

func (r *postgresMessageRepository) GetPublic(ctx context.Context, limit int, before *time.Time) ([]models.PublicMessage, error) {
	rows, err := r.q.GetPublicMessages(ctx, &sqlcgen.GetPublicMessagesParams{
		Before: toPgTimestamp(before),
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, err
	}
	msgs := make([]models.PublicMessage, len(rows))
	for i, row := range rows {
		msgs[i] = models.PublicMessage{
			SenderNickname:    row.SenderNickname,
			RecipientNickname: row.RecipientNickname,
			Ciphertext:        row.Ciphertext,
			CreatedAt:         row.CreatedAt.Time,
		}
	}
	return msgs, nil
}

func toMessage(r *sqlcgen.Message) *models.Message {
	return &models.Message{
		ID:             r.ID,
		ConversationID: r.ConversationID,
		SenderID:       r.SenderID,
		Ciphertext:     r.Ciphertext,
		Nonce:          r.Nonce,
		CreatedAt:      r.CreatedAt.Time,
	}
}

func toPgTimestamp(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}
