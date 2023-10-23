package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"people-food-service/iternal/config"
	repeatable "people-food-service/pkg/utils"
	"time"
)

const (
	query = `CREATE TABLE IF NOT EXISTS person(
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              name VARCHAR(100) NOT NULL ,
                              family_name VARCHAR(100) NOT NULL
);

			CREATE TABLE IF NOT EXISTS food(
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            name VARCHAR(100) NOT NULL,
                            price FLOAT NOT NULL
);

			CREATE TABLE IF NOT EXISTS person_food(
                                   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                   person_id UUID NOT NULL,
                                   food_id UUID NOT NULL,

                                   CONSTRAINT person_fk FOREIGN KEY (person_id) REFERENCES public.person(id) ON DELETE CASCADE ,
                                   CONSTRAINT food_id FOREIGN KEY (food_id) REFERENCES public.food(id),
                                   CONSTRAINT person_food_unique UNIQUE(person_id, food_id)

);`
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int, sc config.StorageConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			log.Println("failed to connect to postgresql")
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)
	if err != nil {
		log.Fatal("err do with tries postgresql")
	}
	//TODO какаято ошибка в этой части, приложение не запускается в докере из-за этого
	//creating db if not exists
	//tx, err := pool.Begin(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//_, err = tx.Exec(ctx, query)
	//if err != nil {
	//	err := tx.Rollback(ctx)
	//	if err != nil {
	//		return nil, err
	//	}
	//	log.Fatal(err)
	//	return nil, err
	//}
	return pool, err
}
