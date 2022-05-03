CREATE TABLE IF NOT EXISTS public.users
(
    id smallint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 32767 CACHE 1 ),
    username character varying(200) NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);