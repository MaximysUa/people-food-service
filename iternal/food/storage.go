package food

import "context"

type Repository interface {
	Create(ctx context.Context, f Food) (string, error)
	FindAll(ctx context.Context) ([]Food, error)
	FindOne(ctx context.Context, name string) (Food, error)
	Update(ctx context.Context, f Food) error
	Delete(ctx context.Context, f Food) error
}
