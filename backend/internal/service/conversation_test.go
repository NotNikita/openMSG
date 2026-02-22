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

func TestConversationService_GetOrCreate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("sorts ids and creates conversation", func(t *testing.T) {
		repo := mocks.NewConversationRepositoryMock(t)
		// "b" > "a", so they get swapped: ("a", "b")
		repo.GetOrCreateMock.Expect(ctx, "a", "b").
			Return(&models.Conversation{
				ID:        "conv-1",
				UserAID:   "a",
				UserBID:   "b",
				CreatedAt: now,
			}, nil)

		svc := service.NewConversationService(repo)
		conv, err := svc.GetOrCreateConversation(ctx, "b", "a") // intentionally swapped

		require.NoError(t, err)
		assert.Equal(t, "conv-1", conv.ID)
	})

	t.Run("same id returns error", func(t *testing.T) {
		repo := mocks.NewConversationRepositoryMock(t)
		svc := service.NewConversationService(repo)

		_, err := svc.GetOrCreateConversation(ctx, "same", "same")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "yourself")
	})

	t.Run("empty ids return error", func(t *testing.T) {
		repo := mocks.NewConversationRepositoryMock(t)
		svc := service.NewConversationService(repo)

		_, err := svc.GetOrCreateConversation(ctx, "", "b")
		require.Error(t, err)
	})
}

func TestConversationService_ListConversations(t *testing.T) {
	ctx := context.Background()

	t.Run("returns conversations for user", func(t *testing.T) {
		repo := mocks.NewConversationRepositoryMock(t)
		repo.GetByUserIDMock.Expect(ctx, "user-1").
			Return([]models.Conversation{{ID: "conv-1"}}, nil)

		svc := service.NewConversationService(repo)
		convs, err := svc.ListConversations(ctx, "user-1")

		require.NoError(t, err)
		assert.Len(t, convs, 1)
	})

	t.Run("empty userID returns error", func(t *testing.T) {
		repo := mocks.NewConversationRepositoryMock(t)
		svc := service.NewConversationService(repo)

		_, err := svc.ListConversations(ctx, "")
		require.Error(t, err)
	})
}
