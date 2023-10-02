package db

import (
	"context"
	"errors"
	"fmt"
	"people-food-service/iternal/food"
	logging "people-food-service/pkg/client/logger"
	"people-food-service/pkg/client/postgresql"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, f food.Food) (string, error) {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		r.logger.Errorf("faild to create a transaction. err: %v\n", err)
		return "", err
	}
	if f.UUID == "" {
		q := `
			INSERT INTO public.food(name, price)
			SELECT $1, $2
			WHERE NOT EXISTS(select name from food where name = $1::varchar)
			RETURNING id
			`
		newFood := r.client.QueryRow(ctx, q, f.Name, f.Price)
		err := newFood.Scan(&f.UUID)
		if err != nil {
			err := r.client.QueryRow(ctx, "SELECT id FROM public.food WHERE name = $1 AND price = $2",
				f.Name, f.Price).Scan(&f.UUID)
			if err != nil {
				r.logger.Errorf("faild to create new food. err: %v\n", err)
				err := tx.Rollback(ctx)
				if err != nil {
					return "", err
				}
				return "", err
			}
			err = tx.Rollback(ctx)
			if err != nil {
				return "", err
			}
			return f.UUID, errors.New("food is already exist")
		}
	} else {
		q := `
			INSERT INTO public.food(id, name, price)
			SELECT $1, $2, $3
			WHERE NOT EXISTS(select name from food where name = $1::varchar)
			RETURNING id
			`
		_, err := r.client.Exec(ctx, q, f.UUID, f.Name, f.Price)
		if err != nil {
			r.logger.Errorf("faild to create new food. err: %v\n", err)
			err := tx.Rollback(ctx)
			if err != nil {
				return "", err
			}
			return "", err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {

		r.logger.Errorf("faild to commit a transaction. err: %v\n", err)
		return "", err
	}
	return f.UUID, nil
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
	var f food.Food
	q := `
		SELECT id, name, price
		FROM public.food
		WHERE name = $1
		`
	row := r.client.QueryRow(ctx, q, name)
	err := row.Scan(&f.UUID, &f.Name, &f.Price)
	if err != nil {
		r.logger.Errorf("Faild with finding food. query:%s\n", formatQuery(q))
		return food.Food{}, err
	}
	return f, nil
}

func (r *repository) Update(ctx context.Context, f food.Food) error {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		r.logger.Errorf("failed to create a transaction, err: %v\n", err)
		return err
	}
	q := `
		UPDATE food
		SET name = $2, price = $3
		WHERE id = $1
		`
	exec, err := r.client.Exec(ctx, q, f.UUID, f.Name, f.Price)
	if err != nil {
		r.logger.Errorf("Failed with exec the query: %s with id: %s\n", formatQuery(q), f.UUID)
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	if exec.RowsAffected() == 0 {
		err = fmt.Errorf("cant find food in table food with id: %s\n", f.UUID)
		r.logger.Errorf(err.Error())
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Errorf("faild to commit a transaction. err: %v\n", err)
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, f food.Food) error {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		r.logger.Errorf("faild to create a transaction. err: %v\n", err)
		return err
	}
	q := `
		DELETE FROM food f 
		WHERE f.id = $1 
		`
	exec, err := r.client.Exec(ctx, q, f.UUID)
	if err != nil {
		r.logger.Errorf("Failed with exec the query: %s with id: %s\n", formatQuery(q), f.UUID)
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	if exec.RowsAffected() == 0 {
		err = fmt.Errorf("cant find food in table food with id: %s\n", f.UUID)
		r.logger.Errorf(err.Error())
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}

		return err
	}
	err = tx.Commit(ctx)
	if err != nil {

		r.logger.Errorf("faild to commit a transaction. err: %v\n", err)
		return err
	}
	return nil
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
