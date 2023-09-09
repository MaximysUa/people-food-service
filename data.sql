CREATE table public.user(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL ,
    family_name VARCHAR(100) NOT NULL,
    CONSTRAINT fk_food FOREIGN KEY (id) REFERENCES public.food(id)


);

CREATE table public.food(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100)
);