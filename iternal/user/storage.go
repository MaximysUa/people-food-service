package user

import "context"

type Repository interface {
	Create(ctx context.Context, user User) error
	FindAll(ctx context.Context) ([]User, error)
	FindOne(ctx context.Context, name, familyName string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, name, familyName string) error
}
