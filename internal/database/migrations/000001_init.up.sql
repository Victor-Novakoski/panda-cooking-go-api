-- Esquema inicial. IDs de usuário e receita são UUID (não dá para adivinhar
-- o próximo); o resto é sequencial.

CREATE TABLE users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT        NOT NULL,
    email         TEXT        NOT NULL UNIQUE CHECK (email = lower(email)),
    password_hash TEXT        NOT NULL,
    image_profile TEXT        NOT NULL DEFAULT '',
    is_adm        BOOLEAN     NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE categories (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT      NOT NULL UNIQUE
);

CREATE TABLE recipes (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL,
    time        TEXT        NOT NULL,
    portions    INTEGER     NOT NULL CHECK (portions > 0),
    user_id     UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    category_id BIGINT      NOT NULL REFERENCES categories (id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- a listagem é das mais novas para as mais antigas
CREATE INDEX recipes_created_at_idx ON recipes (created_at DESC, id DESC);
CREATE INDEX recipes_user_id_idx ON recipes (user_id);
CREATE INDEX recipes_category_id_idx ON recipes (category_id);

CREATE TABLE image_recipes (
    id        BIGSERIAL PRIMARY KEY,
    url       TEXT      NOT NULL,
    recipe_id UUID      NOT NULL REFERENCES recipes (id) ON DELETE CASCADE
);

CREATE INDEX image_recipes_recipe_id_idx ON image_recipes (recipe_id);

CREATE TABLE ingredients (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT      NOT NULL UNIQUE
);

CREATE TABLE ingredient_recipes (
    id            BIGSERIAL PRIMARY KEY,
    amount        TEXT      NOT NULL,
    recipe_id     UUID      NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    ingredient_id BIGINT    NOT NULL REFERENCES ingredients (id)
);

CREATE INDEX ingredient_recipes_recipe_id_idx ON ingredient_recipes (recipe_id);
CREATE INDEX ingredient_recipes_ingredient_id_idx ON ingredient_recipes (ingredient_id);

CREATE TABLE preparations (
    id          BIGSERIAL PRIMARY KEY,
    description TEXT      NOT NULL,
    recipe_id   UUID      NOT NULL REFERENCES recipes (id) ON DELETE CASCADE
);

CREATE INDEX preparations_recipe_id_idx ON preparations (recipe_id);

CREATE TABLE comments (
    id          BIGSERIAL   PRIMARY KEY,
    description TEXT        NOT NULL,
    user_id     UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    recipe_id   UUID        NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- comentários de uma receita, dos mais novos para os mais antigos
CREATE INDEX comments_recipe_id_created_at_idx ON comments (recipe_id, created_at DESC, id DESC);
CREATE INDEX comments_user_id_idx ON comments (user_id);

CREATE TABLE favorite_recipes (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    recipe_id  UUID        NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- o banco impede favorito duplicado, mesmo com dois cliques ao mesmo tempo
    UNIQUE (user_id, recipe_id)
);

CREATE INDEX favorite_recipes_user_id_created_at_idx ON favorite_recipes (user_id, created_at DESC);
CREATE INDEX favorite_recipes_recipe_id_idx ON favorite_recipes (recipe_id);

-- Categorias fixas: a receita sempre aponta para uma delas.
INSERT INTO categories (name) VALUES
    ('Doces e Sobremesas'),
    ('Salgados'),
    ('Massas'),
    ('Carnes'),
    ('Frangos'),
    ('Peixes e Frutos do Mar'),
    ('Sopas e Caldos'),
    ('Saladas'),
    ('Vegetarianos'),
    ('Lanches e Petiscos'),
    ('Bebidas'),
    ('Pães e Bolos');
