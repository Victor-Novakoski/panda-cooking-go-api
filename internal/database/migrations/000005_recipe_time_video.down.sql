ALTER TABLE recipes ADD COLUMN time TEXT;

UPDATE recipes
SET time = CASE
    WHEN time_minutes < 60 THEN time_minutes || ' minutos'
    WHEN time_minutes = 60 THEN '1 hora'
    WHEN time_minutes % 60 = 0 THEN (time_minutes / 60) || ' horas'
    ELSE (time_minutes / 60) || 'h ' || (time_minutes % 60) || 'min'
END;

DROP INDEX IF EXISTS recipes_time_minutes_idx;

ALTER TABLE recipes
    ALTER COLUMN time SET NOT NULL,
    DROP COLUMN time_minutes,
    DROP COLUMN video_url;
