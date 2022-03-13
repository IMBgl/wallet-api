ALTER TABLE public.categories DROP CONSTRAINT IF EXISTS categories_pk;
ALTER TABLE public.categories DROP CONSTRAINT IF EXISTS categories_fk;

DROP TABLE IF EXISTS public.categories;