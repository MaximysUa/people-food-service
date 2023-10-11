package postgresql

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"people-food-service/iternal/config"
	repeatable "people-food-service/pkg/utils"
	"time"
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
	//"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	//driv := dsn + "?sslmode=disable"
	//m, err := migrate.New("file://db/migrations/", driv)
	//if err != nil {
	//	log.Fatal(err)
	//	return nil, err
	//}
	//err = m.Up()
	//if err != nil {
	//	log.Fatal(err)
	//	return nil, err
	//}
	return pool, err
}
