ALTER TABLE public.transactions DROP CONSTRAINT IF EXISTS transactions_pk;
ALTER TABLE public.transactions DROP CONSTRAINT IF EXISTS transactions_users_fk;
ALTER TABLE public.transactions DROP CONSTRAINT IF EXISTS transactions_wallets_fk;

DROP TABLE IF EXISTS public.transactions;