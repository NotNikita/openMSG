package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"app/internal/mocks"
	"app/internal/models"
	"app/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.CreateMock.Expect(ctx, "alice", "pubkey123", "https://example.com/avatar.png").
			Return(&models.User{
				ID:        "uuid-1",
				Nickname:  "alice",
				PublicKey: "pubkey123",
				Avatar:    "https://example.com/avatar.png",
				CreatedAt: now,
			}, nil)

		svc := service.NewUserService(repo)
		user, err := svc.CreateUser(ctx, "alice", "pubkey123", "https://example.com/avatar.png")

		require.NoError(t, err)
		assert.Equal(t, "uuid-1", user.ID)
		assert.Equal(t, "alice", user.Nickname)
	})

	t.Run("empty nickname returns error", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := service.NewUserService(repo)

		_, err := svc.CreateUser(ctx, "", "pubkey123", "avatar")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nickname")
	})

	t.Run("empty public_key returns error", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := service.NewUserService(repo)

		_, err := svc.CreateUser(ctx, "alice", "", "avatar")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "public_key")
	})

	t.Run("empty avatar returns error", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := service.NewUserService(repo)

		_, err := svc.CreateUser(ctx, "alice", "pubkey123", "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "avatar")
	})

	t.Run("repo error propagates", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.CreateMock.Expect(ctx, "alice", "pubkey123", "avatar").
			Return(nil, errors.New("db error"))

		svc := service.NewUserService(repo)
		_, err := svc.CreateUser(ctx, "alice", "pubkey123", "avatar")
		require.Error(t, err)
	})
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()

	t.Run("returns user", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.GetByIDMock.Expect(ctx, "uuid-1").
			Return(&models.User{ID: "uuid-1", Nickname: "alice"}, nil)

		svc := service.NewUserService(repo)
		user, err := svc.GetUser(ctx, "uuid-1")

		require.NoError(t, err)
		assert.Equal(t, "uuid-1", user.ID)
	})

	t.Run("empty id returns error", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := service.NewUserService(repo)

		_, err := svc.GetUser(ctx, "")
		require.Error(t, err)
	})
}
