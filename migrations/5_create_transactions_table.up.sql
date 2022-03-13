CREATE TABLE public.transactions (
	id uuid NOT NULL,
	amount float4 NOT NULL,
	user_id uuid NOT NULL,
    wallet_id uuid NOT NULL,
	currency varchar NOT NULL,
    comment varchar NULL,
    type varchar NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	CONSTRAINT transactions_pk PRIMARY KEY (id),
	CONSTRAINT transactions_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT transactions_wallets_fk FOREIGN KEY (wallet_id) REFERENCES public.wallets(id)
);
