package seed

import "panda-cooking-go-api/internal/model"

type demoUser struct {
	name, email, image string
	isAdm              bool
}

var demoUsers = []demoUser{
	{name: "Chef Maria Silva", email: "maria@pandacooking.com", image: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=400", isAdm: true},
	{name: "João Cozinheiro", email: "joao@pandacooking.com", image: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=400"},
	{name: "Ana Paula Gourmet", email: "ana@pandacooking.com", image: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=400"},
}

type demoIngredient struct {
	name, amount string
}

type demoRecipe struct {
	recipe       model.Recipe
	author       int    // índice em demoUsers
	category     string // nome de uma das categorias criadas pela migration
	images       []string
	ingredients  []demoIngredient // na ordem em que aparecem na receita
	preparations []string
	comments     []demoComment
}

type demoComment struct {
	author int // índice em demoUsers
	text   string
}

var demoRecipes = []demoRecipe{
	// Brigadeiro
	{
		recipe: model.Recipe{
			Name:        "Brigadeiro Tradicional",
			Description: "O docinho mais amado do Brasil! Cremoso, chocolatudo e perfeito para qualquer festa.",
			Time:        "20 minutos",
			Portions:    30,
		},
		author:   0,
		category: "Doces e Sobremesas",
		images: []string{
			"https://images.unsplash.com/photo-1702982852429-e0d0b27eb990?w=800",
		},
		ingredients: []demoIngredient{
			{"Leite condensado", "1 lata (395g)"},
			{"Manteiga", "1 colher de sopa"},
			{"Chocolate em pó", "4 colheres de sopa"},
			{"Chocolate granulado", "100g para cobrir"},
		},
		preparations: []string{
			"Em uma panela, coloque o leite condensado, a manteiga e o chocolate em pó.",
			"Leve ao fogo médio e mexa sem parar até que o brigadeiro comece a desgrudar do fundo da panela.",
			"Despeje o brigadeiro em um prato untado com manteiga e deixe esfriar completamente.",
			"Com as mãos untadas com manteiga, pegue pequenas porções da massa e enrole em formato de bolinhas.",
			"Passe as bolinhas no chocolate granulado e coloque em forminhas de papel.",
		},
		comments: []demoComment{
			{author: 1, text: "Fiz para o aniversário do meu filho e não sobrou nenhum!"},
			{author: 2, text: "Dica: tire do fogo um pouco antes se for comer de colher."},
		},
	},

	// Coxinha
	{
		recipe: model.Recipe{
			Name:        "Coxinha de Frango",
			Description: "Salgado brasileiro clássico, com massa cremosa e recheio suculento de frango desfiado.",
			Time:        "1h 30min",
			Portions:    25,
		},
		author:   1,
		category: "Salgados",
		images: []string{
			"https://images.unsplash.com/photo-1767974963436-2208b3553561?w=800",
		},
		ingredients: []demoIngredient{
			{"Frango", "500g (peito)"},
			{"Cebola", "1 unidade média"},
			{"Alho", "2 dentes"},
			{"Caldo de galinha", "2 tabletes"},
			{"Leite", "2 xícaras"},
			{"Farinha de trigo", "3 xícaras"},
			{"Manteiga", "2 colheres de sopa"},
			{"Farinha de rosca", "300g para empanar"},
			{"Ovos", "2 unidades"},
			{"Sal", "a gosto"},
			{"Cheiro-verde", "a gosto"},
			{"Óleo", "para fritar"},
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
		comments: []demoComment{
			{author: 0, text: "A massa ficou lisinha, ótima de modelar."},
			{author: 2, text: "Coloquei requeijão no recheio e ficou ainda melhor."},
		},
	},

	// Feijoada
	{
		recipe: model.Recipe{
			Name:        "Feijoada Completa",
			Description: "Prato tradicional brasileiro, perfeito para reunir a família. Rica em sabores e muito suculenta!",
			Time:        "3 horas",
			Portions:    10,
		},
		author:   0,
		category: "Carnes",
		images: []string{
			"https://images.unsplash.com/photo-1734771771447-d943e2b5f4d5?w=800",
		},
		ingredients: []demoIngredient{
			{"Feijão", "1kg (preto)"},
			{"Carne suína", "500g (costela)"},
			{"Carne bovina", "300g (carne seca)"},
			{"Linguiça", "400g (calabresa)"},
			{"Bacon", "200g"},
			{"Cebola", "2 unidades"},
			{"Alho", "6 dentes"},
			{"Louro", "3 folhas"},
			{"Sal", "a gosto"},
			{"Pimenta-do-reino", "a gosto"},
			{"Cheiro-verde", "a gosto"},
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
		comments: []demoComment{
			{author: 1, text: "Receita de domingo aqui em casa. Perfeita!"},
		},
	},

	// Pão de Queijo
	{
		recipe: model.Recipe{
			Name:        "Pão de Queijo Mineiro",
			Description: "Clássico de Minas Gerais! Crocante por fora, macio por dentro e com aquele sabor inconfundível de queijo.",
			Time:        "40 minutos",
			Portions:    40,
		},
		author:   2,
		category: "Pães e Bolos",
		images: []string{
			"https://images.unsplash.com/photo-1559141680-d0bd7bc5af84?w=800",
		},
		ingredients: []demoIngredient{
			{"Polvilho azedo", "500g"},
			{"Leite", "1 xícara"},
			{"Óleo", "1/2 xícara"},
			{"Ovos", "3 unidades"},
			{"Queijo mussarela", "150g ralado"},
			{"Queijo parmesão", "100g ralado"},
			{"Sal", "1 colher de chá"},
		},
		preparations: []string{
			"Ferva o leite com o óleo e o sal.",
			"Despeje sobre o polvilho e mexa bem até formar uma massa grossa.",
			"Deixe amornar e adicione os ovos, um a um, misturando bem.",
			"Acrescente os queijos ralados e misture até formar uma massa homogênea.",
			"Com as mãos untadas com óleo, faça bolinhas e coloque em uma assadeira untada.",
			"Leve ao forno preaquecido a 180°C por cerca de 25-30 minutos, até dourar.",
		},
		comments: []demoComment{
			{author: 0, text: "Congelei metade crua e assei na semana seguinte, ficou igual."},
			{author: 1, text: "Melhor que o da padaria."},
		},
	},

	// Moqueca
	{
		recipe: model.Recipe{
			Name:        "Moqueca Capixaba",
			Description: "Ensopado de peixe com urucum, tomate e cheiro-verde. Leve, saboroso e aromático!",
			Time:        "50 minutos",
			Portions:    6,
		},
		author:   0,
		category: "Peixes e Frutos do Mar",
		images: []string{
			"https://images.unsplash.com/photo-1559847844-5315695dadae?w=800",
		},
		ingredients: []demoIngredient{
			{"Peixe", "1kg (badejo ou robalo em postas)"},
			{"Tomate", "4 unidades"},
			{"Cebola", "2 unidades"},
			{"Pimentão", "1 unidade"},
			{"Cheiro-verde", "1 maço"},
			{"Coentro", "a gosto"},
			{"Alho", "4 dentes"},
			{"Limão", "2 unidades"},
			{"Azeite", "4 colheres de sopa"},
			{"Urucum", "1 colher de sopa (colorau)"},
			{"Sal", "a gosto"},
			{"Pimenta-do-reino", "a gosto"},
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
		comments: []demoComment{
			{author: 2, text: "Usei robalo e ficou maravilhosa."},
		},
	},

	// Bolo de Cenoura
	{
		recipe: model.Recipe{
			Name:        "Bolo de Cenoura com Cobertura de Chocolate",
			Description: "Bolo fofinho e úmido de cenoura com aquela cobertura de chocolate cremosa que derrete na boca!",
			Time:        "1 hora",
			Portions:    12,
		},
		author:   2,
		category: "Doces e Sobremesas",
		images: []string{
			"https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=800",
		},
		ingredients: []demoIngredient{
			{"Cenoura", "3 unidades médias"},
			{"Ovos", "3 unidades"},
			{"Óleo", "1/2 xícara"},
			{"Açúcar", "2 xícaras"},
			{"Farinha de trigo", "2 xícaras"},
			{"Fermento em pó", "1 colher de sopa"},
			// Cobertura
			{"Chocolate em pó", "4 colheres de sopa"},
			{"Manteiga", "3 colheres de sopa"},
			{"Leite", "4 colheres de sopa"},
			{"Açúcar", "1 xícara (cobertura)"},
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
		comments: []demoComment{
			{author: 1, text: "O difícil é esperar a cobertura esfriar."},
		},
	},

	// Estrogonofe
	{
		recipe: model.Recipe{
			Name:        "Estrogonofe de Frango",
			Description: "Cremoso e saboroso, esse prato é sempre um sucesso! Perfeito para almoços especiais.",
			Time:        "40 minutos",
			Portions:    6,
		},
		author:   1,
		category: "Frangos",
		images: []string{
			"https://images.unsplash.com/photo-1644592219048-5c070fd3c91c?w=800",
		},
		ingredients: []demoIngredient{
			{"Frango", "1kg (peito em cubos)"},
			{"Cebola", "1 unidade"},
			{"Alho", "3 dentes"},
			{"Tomate", "2 unidades"},
			{"Molho de tomate", "1 lata"},
			{"Creme de leite", "1 lata"},
			{"Mostarda", "1 colher de sopa"},
			{"Ketchup", "2 colheres de sopa"},
			{"Cogumelos", "1 lata (opcional)"},
			{"Azeite", "2 colheres de sopa"},
			{"Sal", "a gosto"},
			{"Pimenta-do-reino", "a gosto"},
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
		comments: []demoComment{
			{author: 0, text: "Rápido e fácil, virou o jantar da semana."},
		},
	},

	// Açaí na Tigela
	{
		recipe: model.Recipe{
			Name:        "Açaí na Tigela",
			Description: "Refrescante e energético! Bowl de açaí cremoso com granola, banana e mel. Perfeito para o café da manhã ou lanche.",
			Time:        "10 minutos",
			Portions:    2,
		},
		author:   2,
		category: "Bebidas",
		images: []string{
			"https://images.unsplash.com/photo-1627308594190-a057cd4bfac8?w=800",
		},
		ingredients: []demoIngredient{
			{"Polpa de açaí", "400g (congelada)"},
			{"Banana", "2 unidades"},
			{"Guaraná em pó", "1 colher de chá (opcional)"},
			{"Mel", "2 colheres de sopa"},
			{"Granola", "4 colheres de sopa"},
			{"Morango", "6 unidades"},
			{"Banana", "1 unidade (cobertura)"},
			{"Coco ralado", "2 colheres de sopa"},
		},
		preparations: []string{
			"No liquidificador, bata a polpa de açaí com as bananas até formar um creme consistente.",
			"Se preferir mais doce, adicione mel ou guaraná em pó para dar energia.",
			"Despeje o açaí batido em tigelas.",
			"Decore com rodelas de banana, morangos, granola e coco ralado.",
			"Regue com mel e sirva imediatamente.",
		},
		comments: []demoComment{
			{author: 1, text: "Bati com um pouco de leite para ficar mais cremoso."},
		},
	},

	// Lasanha
	{
		recipe: model.Recipe{
			Name:        "Lasanha à Bolonhesa",
			Description: "Lasanha tradicional com molho bolonhesa rico, molho branco cremoso e muito queijo. Conforto em forma de comida!",
			Time:        "1h 30min",
			Portions:    8,
		},
		author:   0,
		category: "Massas",
		images: []string{
			"https://images.unsplash.com/photo-1574894709920-11b28e7367e3?w=800",
		},
		ingredients: []demoIngredient{
			{"Massa de lasanha", "500g"},
			{"Carne bovina", "500g (moída)"},
			{"Molho de tomate", "2 latas"},
			{"Cebola", "1 unidade"},
			{"Alho", "3 dentes"},
			{"Tomate", "2 unidades"},
			{"Leite", "1 litro"},
			{"Farinha de trigo", "3 colheres de sopa"},
			{"Manteiga", "3 colheres de sopa"},
			{"Queijo mussarela", "400g"},
			{"Queijo parmesão", "200g"},
			{"Sal", "a gosto"},
			{"Pimenta-do-reino", "a gosto"},
			{"Orégano", "a gosto"},
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
		comments: []demoComment{
			{author: 2, text: "Deixei montada de véspera e só assei no dia. Deu certo!"},
		},
	},

	// Pudim
	{
		recipe: model.Recipe{
			Name:        "Pudim de Leite Condensado",
			Description: "Clássico brasileiro! Sobremesa cremosa e suave com aquela calda de caramelo irresistível.",
			Time:        "1h 20min",
			Portions:    10,
		},
		author:   1,
		category: "Doces e Sobremesas",
		images: []string{
			"https://images.unsplash.com/photo-1702728109878-c61a98d80491?w=800",
		},
		ingredients: []demoIngredient{
			{"Leite condensado", "1 lata"},
			{"Leite", "2 medidas da lata (use a lata de leite condensado)"},
			{"Ovos", "3 unidades"},
			{"Açúcar", "1 xícara (para a calda)"},
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
		comments: []demoComment{
			{author: 2, text: "Sem furinhos! O banho-maria em fogo baixo faz diferença."},
			{author: 0, text: "Clássico que nunca falha."},
		},
	},

	// Picanha
	{
		recipe: model.Recipe{
			Name:        "Picanha na Chapa",
			Description: "O corte preferido do churrasco brasileiro, selado na chapa com a capa de gordura crocante e o miolo rosado.",
			Time:        "30 minutos",
			Portions:    4,
		},
		author:   1,
		category: "Carnes",
		images: []string{
			"https://images.unsplash.com/photo-1702288824191-3003d32484b1?w=800",
		},
		ingredients: []demoIngredient{
			{"Picanha", "1 peça (cerca de 1,2kg)"},
			{"Sal grosso", "a gosto"},
			{"Alho", "4 dentes"},
			{"Manteiga", "1 colher de sopa"},
			{"Alecrim", "2 ramos"},
		},
		preparations: []string{
			"Tire a picanha da geladeira 30 minutos antes e corte em bifes grossos, a favor da gordura, contra as fibras.",
			"Tempere só com sal grosso dos dois lados.",
			"Aqueça bem a chapa ou frigideira de ferro. Comece selando o lado da gordura até ficar dourado e crocante.",
			"Vire e sele cada lado por 3 a 4 minutos para ponto mal passado ou ao ponto.",
			"No fim, junte a manteiga, o alho amassado e o alecrim e regue a carne com a manteiga derretida.",
			"Deixe descansar 5 minutos antes de fatiar e sirva com farofa e vinagrete.",
		},
		comments: []demoComment{
			{author: 0, text: "O segredo é não mexer na carne enquanto sela. Ficou perfeita!"},
		},
	},

	// Frango assado
	{
		recipe: model.Recipe{
			Name:        "Frango Assado com Batatas",
			Description: "Almoço de domingo: frango inteiro temperado na véspera, pele dourada e batatas assadas no próprio caldo.",
			Time:        "1h 40min",
			Portions:    6,
		},
		author:   2,
		category: "Frangos",
		images: []string{
			"https://images.unsplash.com/photo-1768238907887-023b7ac9f450?w=800",
		},
		ingredients: []demoIngredient{
			{"Frango", "1 inteiro (cerca de 2kg)"},
			{"Batata", "1kg (pequenas)"},
			{"Limão", "2 unidades"},
			{"Alho", "6 dentes"},
			{"Páprica defumada", "1 colher de sopa"},
			{"Azeite", "4 colheres de sopa"},
			{"Alecrim", "3 ramos"},
			{"Sal", "a gosto"},
			{"Pimenta-do-reino", "a gosto"},
		},
		preparations: []string{
			"Misture o suco dos limões, o alho amassado, a páprica, o azeite, sal e pimenta.",
			"Espalhe o tempero por todo o frango, inclusive por baixo da pele, e deixe na geladeira por pelo menos 4 horas.",
			"Corte as batatas ao meio e tempere com azeite, sal e alecrim.",
			"Coloque o frango no centro da assadeira, as batatas em volta, e cubra com papel alumínio.",
			"Asse a 200°C por 50 minutos. Retire o papel e asse por mais 30 a 40 minutos, regando com o caldo, até dourar.",
		},
		comments: []demoComment{
			{author: 1, text: "Temperei na noite anterior e ficou muito suculento."},
			{author: 0, text: "Coloquei cebola junto com as batatas, ficou ótimo."},
		},
	},

	// Nhoque
	{
		recipe: model.Recipe{
			Name:        "Nhoque de Batata ao Sugo",
			Description: "Nhoque caseiro, leve e macio, com molho de tomate fresco, manjericão e mussarela derretida.",
			Time:        "1h 10min",
			Portions:    4,
		},
		author:   0,
		category: "Massas",
		images: []string{
			"https://images.unsplash.com/photo-1651773550746-b76fea3cfb70?w=800",
		},
		ingredients: []demoIngredient{
			{"Batata", "1kg (asterix)"},
			{"Farinha de trigo", "1 1/2 xícara"},
			{"Ovos", "1 unidade"},
			{"Sal", "a gosto"},
			{"Tomate", "6 maduros"},
			{"Alho", "3 dentes"},
			{"Azeite", "3 colheres de sopa"},
			{"Manjericão", "1 maço"},
			{"Queijo mussarela", "150g em cubos"},
		},
		preparations: []string{
			"Cozinhe as batatas com casca até ficarem macias. Descasque ainda quentes e passe no espremedor.",
			"Junte o ovo, o sal e a farinha aos poucos, só até a massa parar de grudar nas mãos.",
			"Faça rolinhos, corte em pedaços de 2 cm e passe os nhoques na farinha.",
			"Para o molho, refogue o alho no azeite, junte os tomates picados e cozinhe por 20 minutos. Acerte o sal e junte o manjericão.",
			"Cozinhe os nhoques em água fervente com sal. Quando subirem à superfície, retire com escumadeira.",
			"Misture ao molho, cubra com a mussarela e sirva assim que derreter.",
		},
		comments: []demoComment{
			{author: 2, text: "Fiz no dia 29 com uma nota embaixo do prato, como manda a tradição!"},
		},
	},

	// Torta de limão
	{
		recipe: model.Recipe{
			Name:        "Torta de Limão",
			Description: "Massa amanteigada, creme azedinho de limão e merengue maçaricado por cima. Refrescante e linda.",
			Time:        "1h 30min",
			Portions:    10,
		},
		author:   2,
		category: "Doces e Sobremesas",
		images: []string{
			"https://images.unsplash.com/photo-1667933240742-956b770cd910?w=800",
		},
		ingredients: []demoIngredient{
			{"Biscoito maisena", "200g"},
			{"Manteiga", "100g derretida"},
			{"Leite condensado", "2 latas"},
			{"Limão", "4 unidades (suco e raspas)"},
			{"Claras", "3 unidades"},
			{"Açúcar", "6 colheres de sopa"},
		},
		preparations: []string{
			"Triture os biscoitos, misture com a manteiga e forre o fundo e a lateral de uma forma de fundo removível.",
			"Asse a massa a 180°C por 10 minutos e deixe esfriar.",
			"Misture o leite condensado com o suco dos limões até engrossar e espalhe sobre a massa.",
			"Bata as claras em neve, junte o açúcar aos poucos e bata até formar um merengue firme e brilhante.",
			"Cubra a torta com o merengue e doure com maçarico ou no forno bem quente por alguns minutos.",
			"Finalize com raspas de limão e leve à geladeira por 3 horas antes de servir.",
		},
		comments: []demoComment{
			{author: 1, text: "A melhor torta de limão que já fiz, o creme fica no ponto certo."},
		},
	},

	// Salada de palmito
	{
		recipe: model.Recipe{
			Name:        "Salada de Palmito",
			Description: "Salada colorida e rápida com palmito, tomate, ovos cozidos e um molho simples de limão e azeite.",
			Time:        "20 minutos",
			Portions:    4,
		},
		author:   0,
		category: "Saladas",
		images: []string{
			"https://images.unsplash.com/photo-1626379907504-327b925f4b79?w=800",
		},
		ingredients: []demoIngredient{
			{"Palmito", "1 vidro (300g)"},
			{"Tomate cereja", "200g"},
			{"Ovos", "3 unidades cozidos"},
			{"Alface", "1 pé"},
			{"Cebola roxa", "1/2 unidade"},
			{"Azeite", "3 colheres de sopa"},
			{"Limão", "1 unidade"},
			{"Sal", "a gosto"},
		},
		preparations: []string{
			"Lave e seque bem as folhas de alface e rasgue com as mãos.",
			"Corte o palmito em rodelas, os tomates ao meio, os ovos em quartos e a cebola em fatias finas.",
			"Misture o azeite, o suco do limão e o sal para o molho.",
			"Monte a salada sobre as folhas e regue com o molho na hora de servir.",
		},
		comments: []demoComment{
			{author: 1, text: "Acompanhamento perfeito para o churrasco."},
		},
	},

	// Curry de grão-de-bico
	{
		recipe: model.Recipe{
			Name:        "Curry de Grão-de-Bico",
			Description: "Prato vegetariano cremoso com leite de coco, especiarias e grão-de-bico, servido com arroz soltinho.",
			Time:        "40 minutos",
			Portions:    4,
		},
		author:   2,
		category: "Vegetarianos",
		images: []string{
			"https://images.unsplash.com/photo-1582576163090-09d3b6f8a969?w=800",
		},
		ingredients: []demoIngredient{
			{"Grão-de-bico", "2 xícaras cozido"},
			{"Leite de coco", "1 vidro (200ml)"},
			{"Tomate", "2 unidades"},
			{"Cebola", "1 unidade"},
			{"Alho", "3 dentes"},
			{"Gengibre", "1 colher de chá ralado"},
			{"Curry em pó", "1 colher de sopa"},
			{"Espinafre", "2 xícaras"},
			{"Sal", "a gosto"},
		},
		preparations: []string{
			"Refogue a cebola, o alho e o gengibre até a cebola ficar transparente.",
			"Junte o curry e mexa por 1 minuto para soltar o aroma.",
			"Adicione os tomates picados e cozinhe até desmancharem.",
			"Junte o grão-de-bico e o leite de coco e cozinhe em fogo baixo por 15 minutos.",
			"Desligue o fogo, misture o espinafre e acerte o sal. Sirva com arroz.",
		},
		comments: []demoComment{
			{author: 0, text: "Nem senti falta de carne. Vai entrar no cardápio da semana."},
		},
	},

	// Bolinho de bacalhau
	{
		recipe: model.Recipe{
			Name:        "Bolinho de Bacalhau",
			Description: "Petisco de boteco: crocante por fora, cremoso por dentro e cheio de bacalhau desfiado.",
			Time:        "1 hora",
			Portions:    30,
		},
		author:   1,
		category: "Lanches e Petiscos",
		images: []string{
			"https://images.unsplash.com/photo-1772198537728-b318ab272ffb?w=800",
		},
		ingredients: []demoIngredient{
			{"Bacalhau", "400g dessalgado"},
			{"Batata", "500g"},
			{"Cebola", "1 unidade"},
			{"Alho", "2 dentes"},
			{"Ovos", "2 unidades"},
			{"Salsinha", "a gosto"},
			{"Pimenta-do-reino", "a gosto"},
			{"Óleo", "para fritar"},
		},
		preparations: []string{
			"Cozinhe o bacalhau por 10 minutos, escorra e desfie bem, tirando pele e espinhas.",
			"Cozinhe as batatas e amasse ainda quentes.",
			"Refogue a cebola e o alho picados e misture com o bacalhau, a batata, os ovos, a salsinha e a pimenta.",
			"Modele os bolinhos com duas colheres.",
			"Frite em óleo quente até dourar e escorra em papel toalha.",
		},
		comments: []demoComment{
			{author: 2, text: "Servi com limão e uma cerveja gelada. Sucesso!"},
		},
	},

	// Caipirinha
	{
		recipe: model.Recipe{
			Name:        "Caipirinha de Limão",
			Description: "O drinque mais famoso do Brasil: limão, açúcar, gelo e cachaça. Simples e perfeito.",
			Time:        "5 minutos",
			Portions:    1,
		},
		author:   1,
		category: "Bebidas",
		images: []string{
			"https://images.unsplash.com/photo-1662560094368-821ba04f6d0c?w=800",
		},
		ingredients: []demoIngredient{
			{"Limão", "1 unidade (taiti)"},
			{"Açúcar", "2 colheres de sopa"},
			{"Cachaça", "60ml"},
			{"Gelo", "a gosto"},
		},
		preparations: []string{
			"Corte o limão em gomos e retire a parte branca do meio, que amarga.",
			"No copo, amasse o limão com o açúcar, sem exagero, só para soltar o suco.",
			"Complete com gelo, junte a cachaça e mexa bem.",
		},
		comments: []demoComment{
			{author: 0, text: "Tirar a parte branca do limão muda tudo."},
		},
	},
}

// demoFavorites são as receitas que cada usuário de demonstração favoritou.
var demoFavorites = []struct {
	user   int    // índice em demoUsers
	recipe string // nome da receita
}{
	{user: 0, recipe: "Pão de Queijo Mineiro"},
	{user: 0, recipe: "Coxinha de Frango"},
	{user: 1, recipe: "Brigadeiro Tradicional"},
	{user: 1, recipe: "Feijoada Completa"},
	{user: 1, recipe: "Lasanha à Bolonhesa"},
	{user: 2, recipe: "Moqueca Capixaba"},
	{user: 2, recipe: "Pudim de Leite Condensado"},
	{user: 0, recipe: "Torta de Limão"},
	{user: 1, recipe: "Picanha na Chapa"},
	{user: 2, recipe: "Curry de Grão-de-Bico"},
}
