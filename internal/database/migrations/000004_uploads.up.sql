-- Fotos enviadas do computador (receitas e perfil). Ficam no próprio banco:
-- o front reduz a foto antes de enviar (menos de 1 MB), e assim a API não
-- precisa de disco nem de outro serviço para guardar arquivo.

CREATE TABLE uploads (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    content_type TEXT        NOT NULL,
    data         BYTEA       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX uploads_user_id_idx ON uploads (user_id);
