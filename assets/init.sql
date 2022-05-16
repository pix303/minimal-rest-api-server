CREATE TABLE IF NOT EXISTS public.users
(
    id smallint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 32767 CACHE 1 ),
    username character varying(200) NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS public.items
(
    id smallint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 32767 CACHE 1 ),
    name character varying(64) NOT NULL,
    description character varying(256) NOT NULL,
    pieces integer,
    CONSTRAINT items_pkey PRIMARY KEY (id)
);


--INSERT INTO public.items(name,description,pieces) VALUES('first item','first item description',5);
--INSERT INTO public.items(name,description,pieces) VALUES('second item','second item description',5);