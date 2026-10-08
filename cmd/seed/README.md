# Seed de receitas 🐼👨‍🍳

Dados de demonstração com receitas brasileiras, para o site não abrir vazio.

## O que é criado

- **3 usuários**, todos com a senha `panda-cooking-demo`
- **10 receitas completas**, com fotos (Unsplash), ingredientes com quantidade, modo de preparo e alguns comentários:
  Brigadeiro Tradicional, Coxinha de Frango, Feijoada Completa, Pão de Queijo Mineiro, Moqueca Capixaba, Bolo de Cenoura com Cobertura de Chocolate, Estrogonofe de Frango, Açaí na Tigela, Lasanha à Bolonhesa e Pudim de Leite Condensado

As 12 categorias (Doces e Sobremesas, Salgados, Massas, Carnes, Frangos, Peixes e Frutos do Mar, Sopas e Caldos, Saladas, Vegetarianos, Lanches e Petiscos, Bebidas, Pães e Bolos) vêm da migration, não do seed.

| Nome | E-mail | Admin |
| --- | --- | --- |
| Chef Maria Silva | `maria@pandacooking.com` | ✅ |
| João Cozinheiro | `joao@pandacooking.com` | ❌ |
| Ana Paula Gourmet | `ana@pandacooking.com` | ❌ |

## Como usar

Com `SEED_DEMO=true` (padrão no `docker compose up`), a API cria esses dados sozinha ao subir, se o banco ainda não tiver usuário. Nada que já foi cadastrado é apagado.

Para apagar tudo e recriar:

```bash
make seed
```

`make seed` **apaga todos os dados** do banco (usuários, sessões, receitas, comentários e favoritos) antes de popular. Use só em desenvolvimento.

## Adicionando receitas

Edite `demoRecipes` em `internal/seed/data.go`. O autor é a posição em `demoUsers` e a categoria é o nome:

```go
{
	recipe: model.Recipe{
		Name:        "Sua Receita",
		Description: "Descrição com pelo menos 10 caracteres",
		Time:        "45 minutos",
		Portions:    4,
	},
	author:   0,
	category: "Doces e Sobremesas",
	images:   []string{"https://images.unsplash.com/photo-xxx"},
	ingredients: []demoIngredient{
		{name: "Farinha de trigo", amount: "1 xícara"},
		{name: "Açúcar", amount: "200g"},
	},
	preparations: []string{"Passo 1", "Passo 2"},
},
```

Os textos seguem os mesmos limites da API, e `internal/seed/data_test.go` confere isso. Fotos de outro host funcionam, mas o front só otimiza as do `images.unsplash.com`.
