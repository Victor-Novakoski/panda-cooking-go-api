-- Dados de demonstração protegidos: a senha das contas do seed é pública
-- (decisão 017), então ninguém edita nem apaga essas contas e as receitas
-- que o seed criou. Receita nova publicada por uma conta de demonstração
-- continua com is_demo = false.
ALTER TABLE users ADD COLUMN is_demo BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE recipes ADD COLUMN is_demo BOOLEAN NOT NULL DEFAULT false;

-- bancos que já têm o seed (e-mails de internal/seed/data.go)
UPDATE users SET is_demo = true
WHERE email IN ('maria@pandacooking.com', 'joao@pandacooking.com', 'ana@pandacooking.com');

UPDATE recipes SET is_demo = true
WHERE user_id IN (SELECT id FROM users WHERE is_demo);
