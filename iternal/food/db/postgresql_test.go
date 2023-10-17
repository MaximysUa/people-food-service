package db

import (
	"context"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/sirupsen/logrus"
	"people-food-service/iternal/food"
	logging "people-food-service/pkg/client/logger"
	"testing"
)

func TestRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`
								INSERT INTO public.food(name, price)
								SELECT $1, $2
								WHERE NOT EXISTS(select name from food where name = $1::varchar)
								RETURNING id`).
		WithArgs("Пицца", 7.45).
		WillReturnRows(pgxmock.NewRows([]))
	mock.ExpectCommit()
	l := logrus.New()
	level, _ := logrus.ParseLevel("trace")
	l.SetLevel(level)
	le := logrus.NewEntry(l)
	logger := logging.Logger{le}
	r := NewRepository(mock, &logger)
	_, err = r.Create(context.TODO(), food.Food{
		UUID:  "",
		Name:  "Пицца",
		Price: 7.45,
	})
	if err != nil {
		t.Errorf("error was not expected while updating: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
