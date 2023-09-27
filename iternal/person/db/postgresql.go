package person

import (
	"context"
	"errors"
	"fmt"
	"people-food-service/iternal/food"
	"people-food-service/iternal/person"
	logging "people-food-service/pkg/client/logger"
	"people-food-service/pkg/client/postgresql"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) Create(ctx context.Context, person person.Person) (string, error) {
	var q string
	if person.UUID == "" {
		q = `
		INSERT INTO public.person(name, family_name)
		SELECT $1, $2
		WHERE NOT EXISTS(select name, family_name from person where name = $1::varchar and family_name = $2::varchar)
		RETURNING id
		`
		newPersUUID := r.client.QueryRow(ctx, q, person.Name, person.FamilyName)
		err := newPersUUID.Scan(&person.UUID)
		if err != nil {

			err := r.client.QueryRow(ctx, "SELECT id FROM public.person WHERE name = $1 AND family_name = $2",
				person.Name, person.FamilyName).Scan(&person.UUID)
			if err != nil {
				r.logger.Errorf("faild to create new person. query:%v\n", err)
				return "", err
			}

			return person.UUID, errors.New("person is already exist")
		}
	} else {
		q = `
		INSERT INTO public.person(id, name, family_name)
		SELECT $1, $2, $3
		WHERE NOT EXISTS(select name, family_name from person where name = $2 and family_name = $3)
		RETURNING id
		`
		_, err := r.client.Exec(ctx, q, person.UUID, person.Name, person.FamilyName)

		if err != nil {
			r.logger.Errorf("faild to create new person. query:%v\n", err)
			return "", err
		}
	}

	sq := `
		INSERT INTO public.person_food(person_id, food_id) 
		VALUES($1, $2)
		`

	for _, f := range person.Food {
		_, err := r.client.Exec(ctx, sq, person.UUID, f.UUID)
		if err != nil {
			r.logger.Errorf("faild to insert person food. query:%s\n", formatQuery(sq))
			return "", err
		}

	}
	return person.UUID, nil
}

func (r *repository) FindAll(ctx context.Context) ([]person.Person, error) {
	q := `
		SELECT id, name, family_name
		FROM public.person
		`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		r.logger.Errorf("Faild with finding all people. query:%s\n", formatQuery(q))
		return nil, err
	}
	people := make([]person.Person, 0)
	for rows.Next() {
		var p person.Person

		err := rows.Scan(&p.UUID, &p.Name, &p.FamilyName)
		if err != nil {
			r.logger.Errorf("Faild with scaning person row. err:%v\n", err)
			return nil, err
		}
		sq := `
			SELECT f.id, f.name, f.price
			FROM public.person_food pf
			JOIN public.food f on pf.food_id = f.id
			WHERE pf.person_id = $1;
			`
		personFood := make([]food.Food, 0)
		var f food.Food
		foodRow, err := r.client.Query(ctx, sq, p.UUID)
		if err != nil {
			r.logger.Errorf("Faild with finding food row. query:%s\n", formatQuery(sq))
			return nil, err
		}
		for foodRow.Next() {
			err = foodRow.Scan(&f.UUID, &f.Name, &f.Price)
			if err != nil {
				r.logger.Errorf("Faild with scaning food row. err:%v\n", err)
				return nil, err
			}
			personFood = append(personFood, f)

		}
		p.Food = personFood
		people = append(people, p)
	}
	return people, nil
}

func (r *repository) FindOne(ctx context.Context, name, familyName string) (person.Person, error) {
	var p person.Person
	q := `
		SELECT id, name, family_name
		FROM public.person
		WHERE name = $1 AND family_name = $2
		`
	row := r.client.QueryRow(ctx, q, name, familyName)
	err := row.Scan(&p.UUID, &p.Name, &p.FamilyName)
	if err != nil {
		r.logger.Errorf("Faild with finding person. query:%s\n", formatQuery(q))
		return person.Person{}, err
	}
	sq := `
		SELECT f.id, f.name, f.price
		FROM public.person_food pf
		JOIN public.food f on pf.food_id = f.id
		WHERE pf.person_id = $1;
		`
	personFood := make([]food.Food, 0)
	rows, err := r.client.Query(ctx, sq, p.UUID)
	if err != nil {
		r.logger.Errorf("Faild with finding food row. query:%s\n", formatQuery(sq))

		return person.Person{}, err
	}
	for rows.Next() {
		var f food.Food
		err := rows.Scan(&f.UUID, &f.Name, &f.Price)
		if err != nil {
			r.logger.Errorf("Faild with scaning food row. err:%v\n", err)

			return person.Person{}, err
		}
		personFood = append(personFood, f)
	}
	p.Food = personFood
	return p, nil
}

func (r *repository) Update(ctx context.Context, person person.Person) error {
	q := `	
		UPDATE person 
		SET name = $2, family_name = $3 
		WHERE id = $1
		`
	qDel := `
		DELETE FROM person_food p 
		WHERE p.person_id = $1
		`
	qIns := `
		INSERT INTO person_food(person_id, food_id) 
		VALUES($1, $2) 
		`
	exec, err := r.client.Exec(ctx, q, person.UUID, person.Name, person.FamilyName)
	if err != nil {
		r.logger.Errorf("Failed with exec the query: %s with id: %s\n", formatQuery(q), person.UUID)
		return err
	}
	if exec.RowsAffected() == 0 {
		err = fmt.Errorf("cant find person in table person with id: %s\n", person.UUID)
		r.logger.Errorf(err.Error())
		return err
	}

	_, err = r.client.Exec(ctx, qDel, person.UUID)
	if err != nil {
		//r.logger.Errorf("Failed with exec the query: %s with id: %s\n", formatQuery(q), person.UUID)

	}
	for _, val := range person.Food {
		_, err := r.client.Exec(ctx, qIns, person.UUID, val.UUID)
		if err != nil {
			//r.logger.Errorf("Failed with exec the query: %s with id: %s\n", formatQuery(q), person.UUID)

		}

	}
	return nil
}

func (r *repository) Delete(ctx context.Context, p person.Person) error {
	q := `
		DELETE FROM person p 
		WHERE p.id = $1 
		`
	sq := `
		DELETE FROM person_food pf
		WHERE pf.person_id = $1 
		`
	exec, err := r.client.Exec(ctx, q, p.UUID)
	if err != nil {
		r.logger.Errorf("Failed with exec the query: %s with id: %s\n", formatQuery(q), p.UUID)
		return err
	}
	execSq, err := r.client.Exec(ctx, sq, p.UUID)
	if err != nil {
		r.logger.Errorf("Failed with exec the query: %s with id: %s\n", formatQuery(sq), p.UUID)
		return err
	}
	if exec.RowsAffected() == 0 {
		err = fmt.Errorf("cant find person in table person with id: %s\n", p.UUID)
		r.logger.Errorf(err.Error())
		return err
	}
	if execSq.RowsAffected() == 0 {
		r.logger.Debugf("cant find person in table person_food with id: %s\n", p.UUID)

	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) person.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
