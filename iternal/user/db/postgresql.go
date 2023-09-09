package user

import (
	"context"
	"people-food-service/iternal/user"
	"people-food-service/pkg/client/postgresql"
)

type repository struct {
	client postgresql.Client
}

func (r *repository) Create(ctx context.Context, user user.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) FindAll(ctx context.Context) ([]user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) FindOne(ctx context.Context, name, familyName string) (user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Update(ctx context.Context, user user.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, name, familyName string) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postgresql.Client) user.Repository {
	return &repository{
		client: client,
	}
}
