package main

import (
	"fmt"
	"log"
	"os"
	"panda-cooking-go-api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Conectar ao banco
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=panda_cooking port=5433 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Falha ao conectar ao banco:", err)
	}

	fmt.Println("🐼 Iniciando seed do banco de dados...")

	// Limpar dados existentes (cuidado em produção!)
	fmt.Println("🧹 Limpando dados antigos...")
	db.Exec("TRUNCATE TABLE favorite_recipes, comments, preparations, ingredient_recipes, image_recipes, recipes, ingredients, categories, users RESTART IDENTITY CASCADE")

	// Seed das categorias
	categories := seedCategories(db)
	fmt.Printf("✅ %d categorias criadas\n", len(categories))

	// Seed dos ingredientes
	ingredients := seedIngredients(db)
	fmt.Printf("✅ %d ingredientes criados\n", len(ingredients))

	// Seed dos usuários
	users := seedUsers(db)
	fmt.Printf("✅ %d usuários criados\n", len(users))

	// Seed das receitas
	recipes := seedRecipes(db, users, categories, ingredients)
	fmt.Printf("✅ %d receitas criadas\n", len(recipes))

	fmt.Println("🎉 Seed concluído com sucesso!")
}

func seedCategories(db *gorm.DB) []model.Category {
	categories := []model.Category{
		{Name: "Doces e Sobremesas"},
		{Name: "Salgados"},
		{Name: "Massas"},
		{Name: "Carnes"},
		{Name: "Frangos"},
		{Name: "Peixes e Frutos do Mar"},
		{Name: "Sopas e Caldos"},
		{Name: "Saladas"},
		{Name: "Vegetarianos"},
		{Name: "Lanches e Petiscos"},
		{Name: "Bebidas"},
		{Name: "Pães e Bolos"},
	}

	for i := range categories {
		db.Create(&categories[i])
	}

	return categories
}

func seedIngredients(db *gorm.DB) map[string]model.Ingredient {
	ingredientNames := []string{
		// Básicos
		"Farinha de trigo", "Açúcar", "Sal", "Leite", "Ovos", "Manteiga", "Óleo",
		"Fermento em pó", "Essência de baunilha", "Polvilho azedo",
		// Vegetais
		"Cebola", "Alho", "Tomate", "Batata", "Cenoura", "Pimentão", "Cheiro-verde",
		"Coentro", "Manjericão", "Alface", "Rúcula",
		// Carnes
		"Frango", "Carne bovina", "Carne suína", "Bacon", "Linguiça", "Peixe",
		// Massas e grãos
		"Macarrão", "Arroz", "Feijão", "Lentilha", "Grão-de-bico", "Massa de lasanha",
		// Laticínios
		"Queijo mussarela", "Queijo parmesão", "Creme de leite", "Requeijão", "Iogurte",
		// Temperos
		"Azeite", "Vinagre", "Pimenta-do-reino", "Colorau", "Cominho", "Orégano",
		"Páprica", "Louro", "Urucum", "Mostarda", "Ketchup",
		// Frutas
		"Banana", "Morango", "Limão",
		// Outros
		"Chocolate em pó", "Leite condensado", "Coco ralado", "Farinha de rosca",
		"Molho de tomate", "Extrato de tomate", "Caldo de galinha", "Chocolate granulado",
		"Cogumelos", "Polpa de açaí", "Guaraná em pó", "Mel", "Granola",
	}

	ingredients := make(map[string]model.Ingredient)
	for _, name := range ingredientNames {
		ingredient := model.Ingredient{Name: name}
		db.Create(&ingredient)
		ingredients[name] = ingredient
	}

	return ingredients
}

func seedUsers(db *gorm.DB) []model.User {
	users := []model.User{
		{
			Name:         "Chef Maria Silva",
			Email:        "maria@pandacooking.com",
			Password:     "$2a$10$YourHashedPasswordHere", // Usar bcrypt em produção
			ImageProfile: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=400",
			IsAdm:        true,
		},
		{
			Name:         "João Cozinheiro",
			Email:        "joao@pandacooking.com",
			Password:     "$2a$10$YourHashedPasswordHere",
			ImageProfile: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=400",
			IsAdm:        false,
		},
		{
			Name:         "Ana Paula Gourmet",
			Email:        "ana@pandacooking.com",
			Password:     "$2a$10$YourHashedPasswordHere",
			ImageProfile: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=400",
			IsAdm:        false,
		},
	}

	for i := range users {
		db.Create(&users[i])
	}

	return users
}

func seedRecipes(db *gorm.DB, users []model.User, categories []model.Category, ingredients map[string]model.Ingredient) []model.Recipe {
	recipes := []struct {
		recipe       model.Recipe
		images       []string
		ingredients  map[string]string // nome: quantidade
		preparations []string
	}{
		// Brigadeiro
		{
			recipe: model.Recipe{
				Name:        "Brigadeiro Tradicional",
				Description: "O docinho mais amado do Brasil! Cremoso, chocolatudo e perfeito para qualquer festa.",
				Time:        "20 minutos",
				Portions:    30,
				UserID:      users[0].ID,
				CategoryID:  1, // Doces
			},
			images: []string{
				"https://images.unsplash.com/photo-1558312657-e4d69a7537f4?w=800",
			},
			ingredients: map[string]string{
				"Leite condensado": "1 lata (395g)",
				"Manteiga":         "1 colher de sopa",
				"Chocolate em pó":  "4 colheres de sopa",
				"Chocolate granulado": "100g para cobrir",
			},
			preparations: []string{
				"Em uma panela, coloque o leite condensado, a manteiga e o chocolate em pó.",
				"Leve ao fogo médio e mexa sem parar até que o brigadeiro comece a desgrudar do fundo da panela.",
				"Despeje o brigadeiro em um prato untado com manteiga e deixe esfriar completamente.",
				"Com as mãos untadas com manteiga, pegue pequenas porções da massa e enrole em formato de bolinhas.",
				"Passe as bolinhas no chocolate granulado e coloque em forminhas de papel.",
			},
		},

		// Coxinha
		{
			recipe: model.Recipe{
				Name:        "Coxinha de Frango",
				Description: "Salgado brasileiro clássico, com massa cremosa e recheio suculento de frango desfiado.",
				Time:        "1h 30min",
				Portions:    25,
				UserID:      users[1].ID,
				CategoryID:  2, // Salgados
			},
			images: []string{
				"https://images.unsplash.com/photo-1621939514649-280e2ee25f60?w=800",
			},
			ingredients: map[string]string{
				"Frango":           "500g (peito)",
				"Cebola":           "1 unidade média",
				"Alho":             "2 dentes",
				"Caldo de galinha": "2 tabletes",
				"Leite":            "2 xícaras",
				"Farinha de trigo": "3 xícaras",
				"Manteiga":         "2 colheres de sopa",
				"Farinha de rosca": "300g para empanar",
				"Ovos":             "2 unidades",
				"Sal":              "a gosto",
				"Cheiro-verde":     "a gosto",
				"Óleo":             "para fritar",
			},
			preparations: []string{
				"Cozinhe o frango em água com sal. Depois desfie bem e reserve.",
				"Refogue a cebola e o alho picados. Adicione o frango desfiado e tempere com sal e cheiro-verde. Reserve.",
				"Para a massa: Ferva o leite com a manteiga e os caldos de galinha.",
				"Adicione a farinha de trigo de uma vez e mexa vigorosamente até formar uma massa homogênea.",
				"Deixe a massa esfriar até poder manusear. Pegue pequenas porções e abra na palma da mão.",
				"Coloque uma colher de recheio no centro e feche modelando no formato de coxinha.",
				"Passe as coxinhas no ovo batido e depois na farinha de rosca.",
				"Frite em óleo quente até dourar. Escorra em papel toalha e sirva quente.",
			},
		},

		// Feijoada
		{
			recipe: model.Recipe{
				Name:        "Feijoada Completa",
				Description: "Prato tradicional brasileiro, perfeito para reunir a família. Rica em sabores e muito suculenta!",
				Time:        "3 horas",
				Portions:    10,
				UserID:      users[0].ID,
				CategoryID:  4, // Carnes
			},
			images: []string{
				"https://images.unsplash.com/photo-1628418735687-547127d8dc5e?w=800",
			},
			ingredients: map[string]string{
				"Feijão":       "1kg (preto)",
				"Carne suína":  "500g (costela)",
				"Carne bovina": "300g (carne seca)",
				"Linguiça":     "400g (calabresa)",
				"Bacon":        "200g",
				"Cebola":       "2 unidades",
				"Alho":         "6 dentes",
				"Louro":        "3 folhas",
				"Sal":          "a gosto",
				"Pimenta-do-reino": "a gosto",
				"Cheiro-verde": "a gosto",
			},
			preparations: []string{
				"Deixe o feijão de molho na véspera. Dessalgue a carne seca também.",
				"Em uma panela grande, refogue a cebola e o alho no óleo.",
				"Adicione todas as carnes e deixe selar.",
				"Acrescente o feijão escorrido e água suficiente para cobrir.",
				"Adicione as folhas de louro e tempere com sal e pimenta.",
				"Cozinhe em fogo médio por cerca de 2 horas, até as carnes ficarem macias.",
				"Retire algumas conchas de feijão, amasse e retorne à panela para engrossar o caldo.",
				"Finalize com cheiro-verde picado e sirva com arroz, farofa, couve refogada e laranja.",
			},
		},

		// Pão de Queijo
		{
			recipe: model.Recipe{
				Name:        "Pão de Queijo Mineiro",
				Description: "Clássico de Minas Gerais! Crocante por fora, macio por dentro e com aquele sabor inconfundível de queijo.",
				Time:        "40 minutos",
				Portions:    40,
				UserID:      users[2].ID,
				CategoryID:  12, // Pães e Bolos
			},
			images: []string{
				"https://images.unsplash.com/photo-1613376023733-0a73315d9b06?w=800",
			},
			ingredients: map[string]string{
				"Polvilho azedo":    "500g",
				"Leite":             "1 xícara",
				"Óleo":              "1/2 xícara",
				"Ovos":              "3 unidades",
				"Queijo mussarela":  "150g ralado",
				"Queijo parmesão":   "100g ralado",
				"Sal":               "1 colher de chá",
			},
			preparations: []string{
				"Ferva o leite com o óleo e o sal.",
				"Despeje sobre o polvilho e mexa bem até formar uma massa grossa.",
				"Deixe amornar e adicione os ovos, um a um, misturando bem.",
				"Acrescente os queijos ralados e misture até formar uma massa homogênea.",
				"Com as mãos untadas com óleo, faça bolinhas e coloque em uma assadeira untada.",
				"Leve ao forno preaquecido a 180°C por cerca de 25-30 minutos, até dourar.",
			},
		},

		// Moqueca
		{
			recipe: model.Recipe{
				Name:        "Moqueca Capixaba",
				Description: "Ensopado de peixe com urucum, tomate e cheiro-verde. Leve, saboroso e aromático!",
				Time:        "50 minutos",
				Portions:    6,
				UserID:      users[0].ID,
				CategoryID:  6, // Peixes
			},
			images: []string{
				"https://images.unsplash.com/photo-1559847844-5315695dadae?w=800",
			},
			ingredients: map[string]string{
				"Peixe":        "1kg (badejo ou robalo em postas)",
				"Tomate":       "4 unidades",
				"Cebola":       "2 unidades",
				"Pimentão":     "1 unidade",
				"Cheiro-verde": "1 maço",
				"Coentro":      "a gosto",
				"Alho":         "4 dentes",
				"Limão":        "2 unidades",
				"Azeite":       "4 colheres de sopa",
				"Urucum":       "1 colher de sopa (colorau)",
				"Sal":          "a gosto",
				"Pimenta-do-reino": "a gosto",
			},
			preparations: []string{
				"Tempere o peixe com sal, pimenta, alho amassado e limão. Deixe marinar por 30 minutos.",
				"Em uma panela de barro (ou panela comum), faça camadas de cebola e tomate em rodelas.",
				"Coloque as postas de peixe por cima.",
				"Adicione o pimentão em tiras, o cheiro-verde e o coentro.",
				"Regue com o azeite e polvilhe o colorau.",
				"Tampe a panela e cozinhe em fogo baixo por cerca de 20-25 minutos, sem mexer.",
				"Mexa a panela delicadamente de vez em quando para não grudar.",
				"Sirva com arroz branco e pirão.",
			},
		},

		// Bolo de Cenoura
		{
			recipe: model.Recipe{
				Name:        "Bolo de Cenoura com Cobertura de Chocolate",
				Description: "Bolo fofinho e úmido de cenoura com aquela cobertura de chocolate cremosa que derrete na boca!",
				Time:        "1 hora",
				Portions:    12,
				UserID:      users[2].ID,
				CategoryID:  1, // Doces
			},
			images: []string{
				"https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=800",
			},
			ingredients: map[string]string{
				"Cenoura":          "3 unidades médias",
				"Ovos":             "3 unidades",
				"Óleo":             "1/2 xícara",
				"Açúcar":           "2 xícaras",
				"Farinha de trigo": "2 xícaras",
				"Fermento em pó":   "1 colher de sopa",
				// Cobertura
				"Chocolate em pó":  "4 colheres de sopa",
				"Manteiga":         "3 colheres de sopa",
				"Leite":            "4 colheres de sopa",
				"Açúcar (cobertura)": "1 xícara",
			},
			preparations: []string{
				"Bata no liquidificador a cenoura picada, os ovos e o óleo até ficar homogêneo.",
				"Em uma tigela, misture a farinha e o açúcar.",
				"Adicione a mistura do liquidificador e mexa bem.",
				"Por último, adicione o fermento e misture delicadamente.",
				"Despeje em uma forma untada e enfarinhada.",
				"Asse em forno preaquecido a 180°C por cerca de 40 minutos.",
				"Para a cobertura: em uma panela, misture todos os ingredientes e leve ao fogo mexendo até engrossar.",
				"Despeje a cobertura ainda quente sobre o bolo.",
			},
		},

		// Estrogonofe
		{
			recipe: model.Recipe{
				Name:        "Estrogonofe de Frango",
				Description: "Cremoso e saboroso, esse prato é sempre um sucesso! Perfeito para almoços especiais.",
				Time:        "40 minutos",
				Portions:    6,
				UserID:      users[1].ID,
				CategoryID:  5, // Frangos
			},
			images: []string{
				"https://images.unsplash.com/photo-1604908176997-125f25cc6f3d?w=800",
			},
			ingredients: map[string]string{
				"Frango":           "1kg (peito em cubos)",
				"Cebola":           "1 unidade",
				"Alho":             "3 dentes",
				"Tomate":           "2 unidades",
				"Molho de tomate":  "1 lata",
				"Creme de leite":   "1 lata",
				"Mostarda":         "1 colher de sopa",
				"Ketchup":          "2 colheres de sopa",
				"Cogumelos":        "1 lata (opcional)",
				"Azeite":           "2 colheres de sopa",
				"Sal":              "a gosto",
				"Pimenta-do-reino": "a gosto",
			},
			preparations: []string{
				"Tempere o frango com sal e pimenta.",
				"Em uma panela, aqueça o azeite e refogue a cebola e o alho.",
				"Adicione o frango e deixe dourar.",
				"Acrescente o tomate picado e o molho de tomate. Cozinhe por 10 minutos.",
				"Adicione a mostarda e o ketchup. Misture bem.",
				"Se usar cogumelos, adicione neste momento.",
				"Por último, adicione o creme de leite e mexa até aquecer bem (não deixe ferver).",
				"Sirva com arroz branco e batata palha.",
			},
		},

		// Açaí na Tigela
		{
			recipe: model.Recipe{
				Name:        "Açaí na Tigela",
				Description: "Refrescante e energético! Bowl de açaí cremoso com granola, banana e mel. Perfeito para o café da manhã ou lanche.",
				Time:        "10 minutos",
				Portions:    2,
				UserID:      users[2].ID,
				CategoryID:  11, // Bebidas
			},
			images: []string{
				"https://images.unsplash.com/photo-1590301157890-4810ed352733?w=800",
			},
			ingredients: map[string]string{
				"Polpa de açaí":    "400g (congelada)",
				"Banana":           "2 unidades",
				"Guaraná em pó":    "1 colher de chá (opcional)",
				"Mel":              "2 colheres de sopa",
				"Granola":          "4 colheres de sopa",
				"Morango":          "6 unidades",
				"Banana (cobertura)": "1 unidade",
				"Coco ralado":      "2 colheres de sopa",
			},
			preparations: []string{
				"No liquidificador, bata a polpa de açaí com as bananas até formar um creme consistente.",
				"Se preferir mais doce, adicione mel ou guaraná em pó para dar energia.",
				"Despeje o açaí batido em tigelas.",
				"Decore com rodelas de banana, morangos, granola e coco ralado.",
				"Regue com mel e sirva imediatamente.",
			},
		},

		// Lasanha
		{
			recipe: model.Recipe{
				Name:        "Lasanha à Bolonhesa",
				Description: "Lasanha tradicional com molho bolonhesa rico, molho branco cremoso e muito queijo. Conforto em forma de comida!",
				Time:        "1h 30min",
				Portions:    8,
				UserID:      users[0].ID,
				CategoryID:  3, // Massas
			},
			images: []string{
				"https://images.unsplash.com/photo-1574894709920-11b28e7367e3?w=800",
			},
			ingredients: map[string]string{
				"Massa de lasanha":  "500g",
				"Carne bovina":      "500g (moída)",
				"Molho de tomate":   "2 latas",
				"Cebola":            "1 unidade",
				"Alho":              "3 dentes",
				"Tomate":            "2 unidades",
				"Leite":             "1 litro",
				"Farinha de trigo":  "3 colheres de sopa",
				"Manteiga":          "3 colheres de sopa",
				"Queijo mussarela":  "400g",
				"Queijo parmesão":   "200g",
				"Sal":               "a gosto",
				"Pimenta-do-reino":  "a gosto",
				"Orégano":           "a gosto",
			},
			preparations: []string{
				"Prepare o molho bolonhesa: refogue cebola e alho, adicione a carne moída e deixe dourar.",
				"Acrescente o tomate picado, molho de tomate, sal, pimenta e orégano. Cozinhe por 15 minutos.",
				"Para o molho branco: derreta a manteiga, adicione a farinha e mexa. Acrescente o leite aos poucos mexendo sempre até engrossar.",
				"Cozinhe a massa de lasanha conforme instruções da embalagem.",
				"Em um refratário, monte camadas: molho bolonhesa, massa, molho branco e queijos.",
				"Repita as camadas até acabarem os ingredientes, finalizando com queijo.",
				"Leve ao forno preaquecido a 180°C por 30-40 minutos, até gratinar.",
				"Deixe descansar por 10 minutos antes de servir.",
			},
		},

		// Pudim
		{
			recipe: model.Recipe{
				Name:        "Pudim de Leite Condensado",
				Description: "Clássico brasileiro! Sobremesa cremosa e suave com aquela calda de caramelo irresistível.",
				Time:        "1h 20min",
				Portions:    10,
				UserID:      users[1].ID,
				CategoryID:  1, // Doces
			},
			images: []string{
				"https://images.unsplash.com/photo-1624353365286-3f8d62daad51?w=800",
			},
			ingredients: map[string]string{
				"Leite condensado": "1 lata",
				"Leite":            "2 medidas da lata (use a lata de leite condensado)",
				"Ovos":             "3 unidades",
				"Açúcar":           "1 xícara (para a calda)",
			},
			preparations: []string{
				"Prepare a calda: em uma forma de pudim, coloque o açúcar e leve ao fogo médio até derreter e caramelizar.",
				"Espalhe o caramelo por toda a forma e reserve.",
				"No liquidificador, bata o leite condensado, o leite e os ovos até ficar homogêneo.",
				"Despeje a mistura na forma caramelizada.",
				"Cubra a forma com papel alumínio e leve ao forno em banho-maria a 180°C por cerca de 1 hora.",
				"O pudim está pronto quando espetado com um palito, ele sair limpo.",
				"Deixe esfriar completamente e leve à geladeira por no mínimo 4 horas.",
				"Desenforme em um prato com borda para acomodar a calda.",
			},
		},
	}

	var createdRecipes []model.Recipe

	for _, r := range recipes {
		// Criar a receita
		db.Create(&r.recipe)

		// Adicionar imagens
		for _, imgURL := range r.images {
			image := model.ImageRecipe{
				URL:      imgURL,
				RecipeID: r.recipe.ID,
			}
			db.Create(&image)
		}

		// Adicionar ingredientes
		for ingredientName, amount := range r.ingredients {
			ingredient, exists := ingredients[ingredientName]
			if !exists {
				// Se o ingrediente não existir, criar
				ingredient = model.Ingredient{Name: ingredientName}
				db.Create(&ingredient)
				ingredients[ingredientName] = ingredient
			}

			ingredientRecipe := model.IngredientRecipe{
				Amount:       amount,
				RecipeID:     r.recipe.ID,
				IngredientID: ingredient.ID,
			}
			db.Create(&ingredientRecipe)
		}

		// Adicionar preparações
		for _, prepDesc := range r.preparations {
			preparation := model.Preparation{
				Description: prepDesc,
				RecipeID:    r.recipe.ID,
			}
			db.Create(&preparation)
		}

		createdRecipes = append(createdRecipes, r.recipe)
	}

	return createdRecipes
}
