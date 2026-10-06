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

### 1. Configure a conexão com o banco

Por padrão, o script tenta conectar em:
```
host=localhost user=postgres password=postgres dbname=panda_cooking port=5433
```

Para usar outra configuração, defina a variável de ambiente:
```bash
export DATABASE_URL="host=localhost user=seu_user password=sua_senha dbname=panda_cooking port=5433 sslmode=disable"
```

### 2. Execute o seed

```bash
make seed
```

Ou diretamente:
```bash
go run cmd/seed/main.go
```

## ⚠️ Atenção

O script **limpa todos os dados existentes** antes de popular. Use apenas em ambiente de desenvolvimento!

```sql
TRUNCATE TABLE favorite_recipes, comments, preparations, ingredient_recipes, 
image_recipes, recipes, ingredients, categories, users RESTART IDENTITY CASCADE
```

## 🖼️ Imagens

As imagens são URLs do Unsplash (gratuitas e de alta qualidade). Se quiser usar imagens locais:

1. Baixe as imagens
2. Salve em uma pasta pública (ex: `public/images/recipes/`)
3. Atualize as URLs no arquivo `cmd/seed/main.go`

## 🔑 Usuários de exemplo

| Nome | Email | Admin |
|------|-------|-------|
| Chef Maria Silva | maria@pandacooking.com | ✅ |
| João Cozinheiro | joao@pandacooking.com | ❌ |
| Ana Paula Gourmet | ana@pandacooking.com | ❌ |

**Senha**: As senhas estão hasheadas com bcrypt. Você precisa implementar o hash adequado no código.

## 🛠️ Personalizando

Para adicionar mais receitas, edite o slice `recipes` na função `seedRecipes()` no arquivo `cmd/seed/main.go`.

Exemplo:
```go
{
    recipe: model.Recipe{
        Name:        "Sua Receita",
        Description: "Descrição deliciosa",
        Time:        "45 minutos",
        Portions:    4,
        UserID:      users[0].ID,
        CategoryID:  1,
    },
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
