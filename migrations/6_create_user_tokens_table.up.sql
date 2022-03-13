CREATE TABLE public.user_tokens (
                    id uuid NOT NULL,
                    user_id uuid NOT NULL,
                    hash varchar NOT NULL,
                    expires_at timestamp NOT NULL,
                    created_at timestamp NOT NULL,
                    updated_at timestamp NOT NULL,
                    CONSTRAINT user_tokens_pk PRIMARY KEY (id)
);


ALTER TABLE public.user_tokens ADD CONSTRAINT user_tokens_fk FOREIGN KEY (user_id) REFERENCES public.users(id);