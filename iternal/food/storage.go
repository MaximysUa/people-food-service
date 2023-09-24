package food

import "context"

type Repository interface {
	Create(ctx context.Context, food Food) error
	FindAll(ctx context.Context) ([]Food, error)
	FindOne(ctx context.Context, name string) (Food, error)
	Update(ctx context.Context, food Food) error
	Delete(ctx context.Context, food Food) error
}
