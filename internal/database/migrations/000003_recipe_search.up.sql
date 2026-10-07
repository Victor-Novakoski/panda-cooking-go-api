-- Busca por nome e descrição que ignora acento e maiúscula ("acai" acha
-- "Açaí") e usa índice de trigramas para o LIKE '%termo%'.

CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- unaccent() é STABLE (depende do dicionário em uso) e não pode entrar num
-- índice. Com o dicionário fixo o resultado só depende do texto, então esta
-- versão pode ser IMMUTABLE.
CREATE FUNCTION search_normalize(value TEXT) RETURNS TEXT
    LANGUAGE sql IMMUTABLE PARALLEL SAFE STRICT
    RETURN lower(public.unaccent('public.unaccent'::regdictionary, value));

CREATE INDEX recipes_search_idx ON recipes
    USING gin (search_normalize(name || ' ' || description) gin_trgm_ops);
