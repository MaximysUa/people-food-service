package person

import (
	"context"
	"people-food-service/iternal/food"
	"people-food-service/iternal/person"
	"people-food-service/pkg/client/postgresql"
)

type repository struct {
	client postgresql.Client
}

func (r *repository) Create(ctx context.Context, person person.Person) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) FindAll(ctx context.Context) ([]person.Person, error) {
	q := `
		SELECT id, name, family_name
		FROM public.person
		`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	people := make([]person.Person, 0)
	for rows.Next() {
		var p person.Person

		err := rows.Scan(&p.UUID, &p.Name, &p.FamilyName)
		if err != nil {
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
			return nil, err
		}
		for foodRow.Next() {
			err = foodRow.Scan(&f.UUID, &f.Name, &f.Price)
			if err != nil {
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
		return person.Person{}, err
	}
	for rows.Next() {
		var f food.Food
		err := rows.Scan(&f.UUID, &f.Name, &f.Price)
		if err != nil {
			return person.Person{}, err
		}
		personFood = append(personFood, f)
	}
	p.Food = personFood
	return p, nil
}

func (r *repository) Update(ctx context.Context, person person.Person) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, name, familyName string) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postgresql.Client) person.Repository {
	return &repository{
		client: client,
	}
}
