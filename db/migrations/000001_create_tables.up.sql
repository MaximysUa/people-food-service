BEGIN;
CREATE TABLE IF NOT EXISTS public.person(
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              name VARCHAR(100) NOT NULL ,
                              family_name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS public.food(
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            name VARCHAR(100) NOT NULL,
                            price FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.person_food(
                                   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                   person_id UUID NOT NULL,
                                   food_id UUID NOT NULL,

                                   CONSTRAINT person_fk FOREIGN KEY (person_id) REFERENCES public.person(id) ON DELETE CASCADE ,
                                   CONSTRAINT food_id FOREIGN KEY (food_id) REFERENCES public.food(id),
                                   CONSTRAINT person_food_unique UNIQUE(person_id, food_id)

);
COMMIT;