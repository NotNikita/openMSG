package service_test

import (
	"context"
	"testing"
	"time"

	"app/internal/mocks"
	"app/internal/models"
	"app/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageService_SendMessage(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMessageRepositoryMock(t)
		repo.CreateMock.Expect(ctx, "conv-1", "sender-1", "encrypted", "nonce123").
			Return(&models.Message{
				ID:             "msg-1",
				ConversationID: "conv-1",
				SenderID:       "sender-1",
				Ciphertext:     "encrypted",
				Nonce:          "nonce123",
				CreatedAt:      now,
			}, nil)

		svc := service.NewMessageService(repo)
		msg, err := svc.SendMessage(ctx, "conv-1", "sender-1", "encrypted", "nonce123")

		require.NoError(t, err)
		assert.Equal(t, "msg-1", msg.ID)
		assert.Equal(t, "encrypted", msg.Ciphertext)
	})

	t.Run("empty conversation_id returns error", func(t *testing.T) {
		repo := mocks.NewMessageRepositoryMock(t)
		svc := service.NewMessageService(repo)

		_, err := svc.SendMessage(ctx, "", "sender-1", "encrypted", "nonce")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "conversation_id")
	})

	t.Run("empty sender_id returns error", func(t *testing.T) {
		repo := mocks.NewMessageRepositoryMock(t)
		svc := service.NewMessageService(repo)

		_, err := svc.SendMessage(ctx, "conv-1", "", "encrypted", "nonce")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sender_id")
	})

	t.Run("empty ciphertext returns error", func(t *testing.T) {
		repo := mocks.NewMessageRepositoryMock(t)
		svc := service.NewMessageService(repo)

		_, err := svc.SendMessage(ctx, "conv-1", "sender-1", "", "nonce")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ciphertext")
	})

	t.Run("empty nonce returns error", func(t *testing.T) {
		repo := mocks.NewMessageRepositoryMock(t)
		svc := service.NewMessageService(repo)

		_, err := svc.SendMessage(ctx, "conv-1", "sender-1", "encrypted", "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nonce")
	})
}

func TestMessageService_ListPublicMessages(t *testing.T) {
	ctx := context.Background()

	t.Run("first page (no cursor)", func(t *testing.T) {
		repo := mocks.NewMessageRepositoryMock(t)
		repo.GetPublicMock.Expect(ctx, 100, (*time.Time)(nil)).
			Return([]models.PublicMessage{{SenderNickname: "alice"}}, nil)

		svc := service.NewMessageService(repo)
		msgs, err := svc.ListPublicMessages(ctx, 100, nil)

		require.NoError(t, err)
		assert.Len(t, msgs, 1)
	})

	t.Run("with cursor", func(t *testing.T) {
		cursor := time.Now()
		repo := mocks.NewMessageRepositoryMock(t)
		repo.GetPublicMock.Expect(ctx, 50, &cursor).
			Return([]models.PublicMessage{}, nil)

		svc := service.NewMessageService(repo)
		_, err := svc.ListPublicMessages(ctx, 50, &cursor)
		require.NoError(t, err)
	})

	t.Run("clamps limit above 100 to 100", func(t *testing.T) {
		repo := mocks.NewMessageRepositoryMock(t)
		repo.GetPublicMock.Expect(ctx, 100, (*time.Time)(nil)).Return([]models.PublicMessage{}, nil)

		svc := service.NewMessageService(repo)
		_, err := svc.ListPublicMessages(ctx, 500, nil)
		require.NoError(t, err)
	})
}

func TestMessageService_ListMessages(t *testing.T) {
	ctx := context.Background()

	t.Run("returns messages for conversation", func(t *testing.T) {
		repo := mocks.NewMessageRepositoryMock(t)
		repo.GetByConversationIDMock.Expect(ctx, "conv-1").
			Return([]models.Message{{ID: "msg-1"}}, nil)

		svc := service.NewMessageService(repo)
		msgs, err := svc.ListMessages(ctx, "conv-1")

		require.NoError(t, err)
		assert.Len(t, msgs, 1)
	})

	t.Run("empty conversationID returns error", func(t *testing.T) {
		repo := mocks.NewMessageRepositoryMock(t)
		svc := service.NewMessageService(repo)

		_, err := svc.ListMessages(ctx, "")
		require.Error(t, err)
	})
}
