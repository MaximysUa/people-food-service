package food

import "context"

//go:generate mockgen -source=storage.go -destination=mock/mock.go
type Repository interface {
	Create(ctx context.Context, f Food) (string, error)
	FindAll(ctx context.Context) ([]Food, error)
	FindOne(ctx context.Context, name string) (Food, error)
	Update(ctx context.Context, f Food) error
	Delete(ctx context.Context, f Food) error
}
