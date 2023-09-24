package db

import (
	"context"
	"people-food-service/iternal/food"
	logging "people-food-service/pkg/client/logger"
	"people-food-service/pkg/client/postgresql"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, food food.Food) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) FindAll(ctx context.Context) ([]food.Food, error) {
	q := `
		SELECT id, name, price
		FROM public.food
		`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		r.logger.Errorf("Faild with finding all food. query:%s\n", formatQuery(q))
		return nil, err
	}
	foodList := make([]food.Food, 0)
	for rows.Next() {
		var f food.Food

		err := rows.Scan(&f.UUID, &f.Name, &f.Price)
		if err != nil {
			r.logger.Errorf("Faild with scaning food row. err:%v\n", err)
			return nil, err
		}
		foodList = append(foodList, f)
	}
	return foodList, nil
}

func (r *repository) FindOne(ctx context.Context, name string) (food.Food, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Update(ctx context.Context, food food.Food) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, food food.Food) error {
	//TODO implement me
	panic("implement me")
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) food.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
