CREATE TABLE public.users (
	id uuid NOT NULL,
	"name" varchar NOT NULL,
	"password" varchar NOT NULL,
	email varchar NOT NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL
);