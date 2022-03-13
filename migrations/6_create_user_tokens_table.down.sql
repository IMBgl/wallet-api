ALTER TABLE public.user_tokens DROP CONSTRAINT IF EXISTS user_tokens_pk;
ALTER TABLE public.user_tokens DROP CONSTRAINT IF EXISTS user_tokens_fk;

DROP TABLE IF EXISTS public.user_tokens;