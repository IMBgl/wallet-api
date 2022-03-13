CREATE TABLE public.wallets (
	id uuid NOT NULL,
	"name" varchar NOT NULL,
	user_id uuid NOT NULL,
	currency varchar NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	balance float4 NOT NULL,
	CONSTRAINT wallets_pk PRIMARY KEY (id)
);


-- public.wallets foreign keys

ALTER TABLE public.wallets ADD CONSTRAINT wallets_fk FOREIGN KEY (user_id) REFERENCES public.users(id);