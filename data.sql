DROP table public."person_food";
DROP table public."food";
DROP table public."person";

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

    CONSTRAINT person_fk FOREIGN KEY (person_id) REFERENCES public.person(id) ON DELETE CASCADE ,
    CONSTRAINT food_id FOREIGN KEY (food_id) REFERENCES public.food(id),
    CONSTRAINT person_food_unique UNIQUE(person_id, food_id)

);


DELETE FROM person p WHERE p.name = 'Василий' AND p.family_name = 'Соловьёв';

INSERT INTO food(name, price) VALUES('Пицца', 7.85);
INSERT INTO food(name, price) VALUES('Бурито', 9.55);
INSERT INTO food(name, price) VALUES('Шавуха', 15.22);

INSERT INTO public.person(name, family_name) VALUES('Василий', 'Уткин');
INSERT INTO public.person(name, family_name) VALUES('Василий', 'Соловьёв');
INSERT INTO public.person(name, family_name) VALUES('Игорь', 'Адамов');
--0275a687-a26c-49da-a9ae-8070b3a08abe
--0275a687-a26c-49da-a9ae-8070b3a08abe
--0275a687-a26c-49da-a9ae-8070b3a08abe
INSERT INTO public.person_food(person_id, food_id) VALUES('bf69e26f-88c4-4695-b30c-cb27989acbf1',
                                                          '0a5d676b-e8a0-4ee1-8994-5d8b3d954c21');
INSERT INTO public.person_food(person_id, food_id) VALUES('bf69e26f-88c4-4695-b30c-cb27989acbf1',
                                                          '4b982cf3-9b2d-437b-8be2-cb0c28b2bd7a');
INSERT INTO public.person_food(person_id, food_id) VALUES('bf69e26f-88c4-4695-b30c-cb27989acbf1',
                                                          '4c1e120b-1c18-49fd-add6-d2772aa9f541');
INSERT INTO public.person_food(person_id, food_id) VALUES('215386a0-6332-4f00-83ae-2236ab656862',
                                                          '4c1e120b-1c18-49fd-add6-d2772aa9f541');

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

