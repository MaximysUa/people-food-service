DROP table public."person";
DROP table public."food";

CREATE table public.person(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL ,
    family_name VARCHAR(100) NOT NULL
);

CREATE table public.food(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    price FLOAT NOT NULL
);

CREATE TABLE public.person_food(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id UUID NOT NULL,
    food_id UUID NOT NULL,

    CONSTRAINT person_fk FOREIGN KEY (person_id) REFERENCES public.person(id),
    CONSTRAINT food_id FOREIGN KEY (food_id) REFERENCES public.food(id),
    CONSTRAINT person_food_unique UNIQUE(person_id, food_id)

);


DELETE FROM person p WHERE p.name = 'Игорь' AND p.family_name = 'Адамов';

INSERT INTO food(name, price) VALUES('Пицца', 7.85);
INSERT INTO food(name, price) VALUES('Бурито', 9.55);
INSERT INTO food(name, price) VALUES('Шавуха', 15.22);

INSERT INTO public.person(name, family_name) VALUES('Василий', 'Уткин');
INSERT INTO public.person(name, family_name) VALUES('Василий', 'Соловьёв');
INSERT INTO public.person(name, family_name) VALUES('Игорь', 'Адамов');
--0275a687-a26c-49da-a9ae-8070b3a08abe
--0275a687-a26c-49da-a9ae-8070b3a08abe
--0275a687-a26c-49da-a9ae-8070b3a08abe
INSERT INTO public.person_food(person_id, food_id) VALUES('0275a687-a26c-49da-a9ae-8070b3a08abe',
                                                          'd41b9758-f344-447f-b512-cc35b89c23e9');
INSERT INTO public.person_food(person_id, food_id) VALUES('0275a687-a26c-49da-a9ae-8070b3a08abe',
                                                          '41b72d27-c250-4a3a-8c0b-8a7de570a564');
INSERT INTO public.person_food(person_id, food_id) VALUES('0275a687-a26c-49da-a9ae-8070b3a08abe',
                                                          '63aa08fd-15b1-4cb6-af21-e1e40cbadc6c');
INSERT INTO public.person_food(person_id, food_id) VALUES('d8830326-5f1c-49db-942c-63941b7615b8',
                                                          '63aa08fd-15b1-4cb6-af21-e1e40cbadc6c');

SELECT f.id, f.name, f.price
FROM public.person_food pf
JOIN public.food f on pf.food_id = f.id
WHERE pf.person_id = '0275a687-a26c-49da-a9ae-8070b3a08abe';

SELECT p.id, p.name, p.family_name
FROM public.person_food pf
JOIN public.person p on pf.person_id = p.id
WHERE pf.person_id = '0275a687-a26c-49da-a9ae-8070b3a08abe';

SELECT id, name, family_name
FROM public.person
WHERE name = 'Василий' AND family_name = 'Уткин';

INSERT INTO public.person(name, family_name)
SELECT 'Диман', 'Рек'
WHERE NOT EXISTS(select name, family_name from person where name = 'Диман' and family_name = 'Рек')
RETURNING id

