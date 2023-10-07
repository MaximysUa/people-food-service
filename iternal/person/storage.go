package person

import "context"

//go:generate mockgen -source=storage.go -destination=mock/mock.go
type Repository interface {
	Create(ctx context.Context, person Person) (string, error)
	FindAll(ctx context.Context) ([]Person, error)
	FindOne(ctx context.Context, name, familyName string) (Person, error)
	Update(ctx context.Context, person Person) error
	Delete(ctx context.Context, person Person) error
}
