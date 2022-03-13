ALTER TABLE public.wallets DROP CONSTRAINT IF EXISTS wallets_pk;
ALTER TABLE public.wallets DROP CONSTRAINT IF EXISTS wallets_fk;

DROP TABLE IF EXISTS public.wallets;