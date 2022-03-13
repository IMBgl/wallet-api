CREATE TABLE public.categories (
	id uuid NOT NULL,
	"name" varchar NOT NULL,
	user_id uuid NOT NULL,
	balance float4 NOT NULL,
	currency varchar NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	CONSTRAINT categories_pk PRIMARY KEY (id),
	CONSTRAINT categories_fk FOREIGN KEY (user_id) REFERENCES public.users(id)
);
