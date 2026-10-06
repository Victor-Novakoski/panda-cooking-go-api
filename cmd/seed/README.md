# Seed de Receitas 🐼👨‍🍳

Script para popular o banco de dados com receitas brasileiras reais e deliciosas!

## 📋 O que é criado

- **12 categorias**: Doces, Salgados, Massas, Carnes, Frangos, Peixes, Sopas, Saladas, Vegetarianos, Lanches, Bebidas, Pães e Bolos
- **50+ ingredientes**: Ingredientes comuns da culinária brasileira
- **3 usuários de exemplo**: Chefs com fotos de perfil
- **10 receitas completas** incluindo:
  - Brigadeiro Tradicional
  - Coxinha de Frango
  - Feijoada Completa
  - Pão de Queijo Mineiro
  - Moqueca Capixaba
  - Bolo de Cenoura com Chocolate
  - Estrogonofe de Frango
  - Açaí na Tigela
  - Lasanha à Bolonhesa
  - Pudim de Leite Condensado

Cada receita inclui:
- ✅ Nome e descrição detalhada
- ✅ Tempo de preparo
- ✅ Número de porções
- ✅ Imagens reais (Unsplash)
- ✅ Lista completa de ingredientes com quantidades
- ✅ Modo de preparo passo a passo

## 🚀 Como usar

Com `SEED_DEMO=true` (padrão no `docker compose up`), a API cria esses dados sozinha ao subir, se o banco ainda não tiver usuário.

Para apagar tudo e recriar:

```bash
make seed
```

## ⚠️ Atenção

`make seed` **apaga todos os dados** do banco antes de popular. Use apenas em desenvolvimento.

## 🖼️ Imagens

As imagens são URLs do Unsplash (gratuitas e de alta qualidade). Se quiser usar imagens locais:

1. Baixe as imagens
2. Salve em uma pasta pública (ex: `public/images/recipes/`)
3. Atualize as URLs em `internal/seed/data.go`

## 🔑 Usuários de exemplo

Senha de todos: `panda-cooking-demo`.

| Nome | Email | Admin |
|------|-------|-------|
| Chef Maria Silva | maria@pandacooking.com | ✅ |
| João Cozinheiro | joao@pandacooking.com | ❌ |
| Ana Paula Gourmet | ana@pandacooking.com | ❌ |

## 🛠️ Personalizando

Para adicionar mais receitas, edite `demoRecipes` em `internal/seed/data.go`. O autor é a posição em `demoUsers` e a categoria é o nome.

Exemplo:
```go
{
    recipe: model.Recipe{
        Name:        "Sua Receita",
        Description: "Descrição deliciosa",
        Time:        "45 minutos",
        Portions:    4,
    },
    author:   0,
    category: "Doces e Sobremesas",
    images: []string{
        "https://images.unsplash.com/photo-xxx",
    },
    ingredients: map[string]string{
        "Ingrediente 1": "1 xícara",
        "Ingrediente 2": "200g",
    },
    preparations: []string{
        "Passo 1",
        "Passo 2",
    },
},
```

## 📸 Fontes de imagens gratuitas

- [Unsplash](https://unsplash.com/) - Usado neste seed
- [Pexels](https://www.pexels.com/)
- [Pixabay](https://pixabay.com/)

Busque por: "brazilian food", "coxinha", "brigadeiro", "feijoada", etc.
