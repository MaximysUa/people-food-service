package person

import (
	"context"
	"fmt"
	"log"
	"people-food-service/iternal/food"
	"people-food-service/iternal/person"
	"people-food-service/pkg/client/postgresql"
	"strings"
)

type repository struct {
	client postgresql.Client
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
func (r *repository) Create(ctx context.Context, person person.Person) error {
	q := `
		INSERT INTO public.person(name, family_name) 
		VALUES($1, $2)
		`

	newPers := r.client.QueryRow(ctx, q, person.Name, person.FamilyName)
	err := newPers.Scan(&person.UUID)
	if err != nil {
		log.Printf("faild to scan new person. query:%s\n", formatQuery(q))
		return err
	}
	sq := `
		INSERT INTO public.person_food(person_id, food_id) 
		VALUES($1, $2)
		`

	for _, f := range person.Food {
		r.client.QueryRow(ctx, sq, person.UUID, f.UUID)
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) ([]person.Person, error) {
	q := `
		SELECT id, name, family_name
		FROM public.person
		`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		log.Printf("Faild with finding all people. query:%s\n", formatQuery(q))
		return nil, err
	}
	people := make([]person.Person, 0)
	for rows.Next() {
		var p person.Person

		err := rows.Scan(&p.UUID, &p.Name, &p.FamilyName)
		if err != nil {
			log.Printf("Faild with scaning person row. err:%v\n", err)
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
			log.Printf("Faild with finding food row. query:%s\n", formatQuery(sq))
			return nil, err
		}
		for foodRow.Next() {
			err = foodRow.Scan(&f.UUID, &f.Name, &f.Price)
			if err != nil {
				log.Printf("Faild with scaning food row. err:%v\n", err)
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
		log.Printf("Faild with finding person. query:%s\n", formatQuery(q))
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
		log.Printf("Faild with finding food row. query:%s\n", formatQuery(sq))

		return person.Person{}, err
	}
	for rows.Next() {
		var f food.Food
		err := rows.Scan(&f.UUID, &f.Name, &f.Price)
		if err != nil {
			log.Printf("Faild with scaning food row. err:%v\n", err)

			return person.Person{}, err
		}
		personFood = append(personFood, f)
	}
	p.Food = personFood
	return p, nil
}

// Update TODO придумать как быть с таблицей person_food
func (r *repository) Update(ctx context.Context, person person.Person) error {
	q := `	
		UPDATE person 
		SET name = $2, family_name = $3 
		WHERE id = $1
		`
	//sq := `UPDATE person_food SET name = $2, family_name = $3 WHERE id = $1;`
	exec, err := r.client.Exec(ctx, q, person.UUID, person.Name, person.FamilyName)
	if err != nil {
		log.Printf("Failed with exec the query: %s with id: %s\n", formatQuery(q), person.UUID)
		return err
	}
	if exec.RowsAffected() == 0 {
		err = fmt.Errorf("cant find person in table person with id: %s\n", person.UUID)
		log.Println(err)
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	q := `
		DELETE FROM person p 
		WHERE p.id = $1 
		`
	sq := `
		DELETE FROM person_food pf
		WHERE pf.person_id = $1 
		`
	exec, err := r.client.Exec(ctx, q, id)
	if err != nil {
		log.Printf("Failed with exec the query: %s with id: %s\n", formatQuery(q), id)
		return err
	}
	execSq, err := r.client.Exec(ctx, sq, id)
	if err != nil {
		log.Printf("Failed with exec the query: %s with id: %s\n", formatQuery(sq), id)
		return err
	}
	if exec.RowsAffected() == 0 {
		err = fmt.Errorf("cant find person in table person with id: %s\n", id)
		log.Println(err)
		return err
	}
	if execSq.RowsAffected() == 0 {
		err = fmt.Errorf("cant find person in table person_food with id: %s\n", id)
		log.Println(err)
		return err
	}

	return nil
}

func NewRepository(client postgresql.Client) person.Repository {
	return &repository{
		client: client,
	}
}
