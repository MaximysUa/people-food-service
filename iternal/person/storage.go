package person

import "context"

type Repository interface {
	Create(ctx context.Context, person Person) error
	FindAll(ctx context.Context) ([]Person, error)
	FindOne(ctx context.Context, name, familyName string) (Person, error)
	Update(ctx context.Context, person Person) error
	Delete(ctx context.Context, person Person) error
}
