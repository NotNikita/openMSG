package repository

import (
	"context"

	"app/internal/models"
	"app/internal/repository/sqlcgen"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate minimock -i UserRepository -o ../mocks/user_repo_mock.go -n UserRepositoryMock -p mocks

type UserRepository interface {
	Create(ctx context.Context, nickname, publicKey, avatar string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.User, error)
}

type postgresUserRepository struct {
	q *sqlcgen.Queries
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{q: sqlcgen.New(pool)}
}

func (r *postgresUserRepository) Create(ctx context.Context, nickname, publicKey, avatar string) (*models.User, error) {
	row, err := r.q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		Nickname:  nickname,
		PublicKey: publicKey,
		Avatar:    avatar,
	})
	if err != nil {
		return nil, err
	}
	return toUser(row), nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}
	return toUser(row), nil
}

func (r *postgresUserRepository) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	rows, err := r.q.GetAllUsers(ctx, &sqlcgen.GetAllUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	users := make([]models.User, len(rows))
	for i, row := range rows {
		users[i] = *toUser(row)
	}
	return users, nil
}

func toUser(r *sqlcgen.User) *models.User {
	return &models.User{
		ID:        r.ID,
		Nickname:  r.Nickname,
		PublicKey: r.PublicKey,
		Avatar:    r.Avatar,
		CreatedAt: r.CreatedAt.Time,
	}
}
