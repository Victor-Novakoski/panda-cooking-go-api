-- Vídeo da receita: link do YouTube, vazio quando não tem.
ALTER TABLE recipes ADD COLUMN video_url TEXT NOT NULL DEFAULT '';

-- Tempo de preparo em minutos no lugar do texto livre ("40", "20 minutos",
-- "1h 30min", "3 horas"): o site passa a mostrar sempre do mesmo jeito e a
-- listagem ganha ordenação e filtro por tempo.
ALTER TABLE recipes ADD COLUMN time_minutes INTEGER;

UPDATE recipes
SET time_minutes = (
    SELECT CASE
        -- só o número: minutos
        WHEN t ~ '^\s*\d{1,4}\s*$' THEN t::int
        ELSE
            -- horas ("1h", "3 horas")
            COALESCE((regexp_match(t, '(\d{1,4})\s*h'))[1]::int, 0) * 60
            -- minutos escritos ("30min", "20 minutos") ou soltos depois das horas ("1h30")
            + COALESCE(
                (regexp_match(t, '(\d{1,4})\s*m'))[1]::int,
                (regexp_match(t, 'h[a-z]*\s*e?\s*(\d{1,4})\s*$'))[1]::int,
                0)
    END
    FROM (SELECT lower(time) AS t) AS parsed
);

-- o que não deu para entender vira meia hora; nada passa de dois dias
UPDATE recipes SET time_minutes = 30 WHERE time_minutes IS NULL OR time_minutes < 1;
UPDATE recipes SET time_minutes = 2880 WHERE time_minutes > 2880;

ALTER TABLE recipes
    ALTER COLUMN time_minutes SET NOT NULL,
    ADD CONSTRAINT recipes_time_minutes_check CHECK (time_minutes BETWEEN 1 AND 2880),
    DROP COLUMN time;

-- ordenar pelas mais rápidas e filtrar por tempo máximo
CREATE INDEX recipes_time_minutes_idx ON recipes (time_minutes);
