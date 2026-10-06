# CLAUDE.md

Guia de contexto para novos projetos neste ecossistema Go.

## Stack & Dependências Principais

- **go-swagger** (Swagger 2.0 / OpenAPI v2) para geração de server
- **GORM** + PostgreSQL (`gorm.io/driver/postgres`)
- **go-utils** interno: `github.com/gll-sistemas/go-utils` (auth, db, logger, httperror, etc.)
- **RQL** (`github.com/a8m/rql`) para filtros via query string
- **Atlas** para migrations (`ariga.io/atlas-provider-gorm`)
- **logrus** para logging estruturado
- **testcontainers** + **httpexpect** para testes E2E

## Estrutura de Diretórios

```
.
├── spec/
│   └── swagger-api.yaml          # Spec Swagger 2.0 — editado manualmente
├── generated/swagger/            # Código gerado pelo go-swagger — NÃO editar
│   ├── cmd/api-server/main.go
│   └── restapi/
│       ├── ops/                  # Interfaces e tipos de params/responses gerados
│       ├── dto/                  # DTOs gerados
│       └── server.go
├── internal/restapi/
│   ├── handler/
│   │   ├── handler.go            # Struct Handler + DI + ConfigureServer/Middleware
│   │   └── operation/            # Um arquivo por recurso (company.go, address.go…)
│   ├── service/                  # Lógica de negócio
│   ├── repository/               # Acesso ao banco (GORM)
│   ├── model/                    # Models GORM + tags rql
│   └── config/                   # Config YAML opcional
├── database/migrations/          # SQL gerado pelo Atlas
├── atlas.hcl
└── Makefile
```

## Iniciando um Projeto do Zero

```bash
mkdir minha-api && cd minha-api
go mod init minha-api          # nome simples, sem github.com/...
```

## Geração de Código (go-swagger)

O `--implementation-package` deve bater exatamente com o nome do módulo no `go.mod`.

```bash
swagger generate server \
  --spec=./spec/swagger-api.yaml \
  --keep-spec-order \
  --strict-additional-properties \
  --target=./generated/swagger \
  --default-scheme=http \
  --with-flatten=full \
  --main-package=api-server \
  --api-package=ops \
  --model-package=restapi/dto \
  --server-package=restapi \
  --implementation-package minha-api/internal/restapi/handler \
  --principal=github.com/gll-sistemas/go-utils/auth.AuthUser
```

O Makefile deve criar a pasta antes de gerar:

```makefile
swagger:
	rm -fr ./generated/swagger
	mkdir -p generated/swagger
	swagger generate server ...
```

Após gerar, ver `generated/swagger/restapi/auto_configure_*.go` — ele mostra quais interfaces o handler precisa implementar e qual o nome exato do tipo da API (ex: `ops.MinhaAPIAPI`).

Makefile target: `make swagger`

## Padrão de Camadas

### Handler (`internal/restapi/handler/handler.go`)

Sem banco (API simples):

```go
type Handler struct {
    *operation.RecursoOperations
}

func New() *Handler {
    return &Handler{
        RecursoOperations: operation.NewRecursoOperations(),
    }
}

// métodos Configurable — sempre necessários, podem ficar vazios em APIs simples
func (h *Handler) ConfigureFlags(api *ops.MinhaAPIAPI)                    {}
func (h *Handler) ConfigureTLS(tlsConfig *tls.Config)                     {}
func (h *Handler) ConfigureServer(s *http.Server, scheme, addr string)    {}
func (h *Handler) CustomConfigure(api *ops.MinhaAPIAPI)                   {}
func (h *Handler) SetupMiddlewares(handler http.Handler) http.Handler      { return handler }
func (h *Handler) SetupGlobalMiddleware(handler http.Handler) http.Handler { return handler }
```

Com banco (padrão completo):

```go
type Handler struct {
    *db.Conn
    *auth.Auth
    *logger.Logger

    *operation.CompanyOperations
}

func New() *Handler {
    dbConn := db.NewConn()
    return &Handler{
        Conn:              dbConn,
        Auth:              auth.NewAuth(),
        Logger:            logger.NewLogger(),
        CompanyOperations: operation.NewCompanyOperations(
            service.NewCompanyService(
                repository.NewCompanyRepository(dbConn),
            ),
        ),
    }
}
```

### Operation (`internal/restapi/handler/operation/recurso.go`)

```go
func (o *RecursoOperations) CreateRecurso(params ops.CreateRecursoParams) middleware.Responder {
    // campos obrigatórios chegam como ponteiros — desreferenciar antes de usar
    valor := *params.RequestBody.Campo

    return ops.NewCreateRecursoCreated().WithPayload(result)
}
```

Com autenticação (quando o endpoint tem security no yaml):

```go
func (o *CompanyOperations) CreateCompany(params ops.CreateCompanyParams, principal *auth.AuthUser) middleware.Responder {
    dto, err := o.service.CreateCompany(params.RequestBody, principal)
    if err != nil {
        return httperror.HandleHTTPError(err)
    }
    return ops.NewCreateCompanyCreated().WithPayload(dto)
}
```

### Service → Repository → Model

Cadeia simples sem frameworks: Service recebe DTOs, chama Repository, retorna DTOs ou erros.

## Models GORM

```go
type Company struct {
    db.BaseModel   // CreatedAt, UpdatedAt, DisabledAt (soft delete), TrackingInfo

    ID   string `gorm:"primary_key;type:varchar(18);column:id;<-:create"`
    Name string `gorm:"type:varchar(255);not null;column:name" rql:"filter,sort"`
}
```

- IDs gerados via `db.GenerateOrValidateID(galileuid.GIDTXxx, id)` — formato prefixado tipo `CPY-20250414-abc`
- `rql:"filter,sort"` expõe o campo para filtros via criteria

## Filtros (RQL)

Query param: `?criteria={"filter":{"name":{"$eq":"Acme"}},"limit":10,"offset":0}`

No repository:

```go
filters, err := utils.ParseCriteria(criteria, model.Company{})
dbQuery := r.gorm.GetConn().Where(filters.RQLFilter, filters.RQLArgs...).Find(&companies)
```

## Tratamento de Erros

```go
// No handler, sempre:
return httperror.HandleHTTPError(err)

// No service/repository, retornar erros tipados:
// - gorm.ErrRecordNotFound → 404
// - pgErr.Code "23505"    → 400 (unique violation)
// - "bad request: ..."    → 400 (prefixo convencionado)
// - qualquer outro        → 500 com traceId
```

## Configuração (go-flags)

```go
type myOpts struct {
    PostgresDSN string `long:"postgres-dsn" env:"APP_POSTGRES_DSN" required:"true"`
    LogLevel    string `long:"log-level"    env:"APP_LOG_LEVEL"    default:"info"`
}
```

O go-swagger chama `ConfigureFlags()` no handler para registrar grupos de opções.

## Migrations (Atlas)

```bash
make generate-migrations   # gera diff schema → nova migration SQL
make apply-migrations      # aplica pendentes
```

`atlas.hcl` carrega models GORM via `ariga.io/atlas-provider-gorm`.

## Hot Reload (air)

Criar `.air.toml` na raiz. O `pre_cmd` sobe o docker antes e `post_cmd` derruba ao parar.

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = ".temp"

[build]
  args_bin = [
    "--port=8080",
    "--log-format=text",
    "--secret-key=my-secret-key",
    "--issuer-name=\"NomeDoServico\"",
    "--postgres-dsn='host=localhost dbname=devel user=postgres password=postgres port=5432'",
    "--environment-id=P",
    "--instance-id=n",
  ]
  bin = "./.temp/api-server"
  pre_cmd = ["docker compose -f docker-compose-devel.yml up -d"]
  cmd = "go build -o ./.temp/api-server ./generated/swagger/cmd/api-server"
  post_cmd = ["docker compose -f docker-compose-devel.yml down"]
  delay = 0
  exclude_dir = ["assets", ".temp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

Flags em `args_bin` devem existir no projeto — só adicionar flags do go-utils (`--log-format`, `--postgres-dsn`, etc.) se o handler registrar esses option groups em `ConfigureFlags`. Em APIs simples sem banco, basta `--port`.

## Makefile Targets Padrão

| Target | O que faz |
|---|---|
| `swagger` | Regenera código do spec |
| `dev` | Sobe com hot reload (air) |
| `test` | Roda testes |
| `generate-migrations` | Gera migration SQL via Atlas (projetos com banco) |
| `apply-migrations` | Aplica migrations (projetos com banco) |

## Convenções de Nomenclatura

- **Models**: `Company`, `Address`
- **Services**: `CompanyService`
- **Repositories**: `CompanyRepositoryImpl`
- **Operations**: `CompanyOperations`
- **Métodos CRUD**: `CreateX`, `GetXs`, `GetXByID`, `UpdateXByID`, `DeleteXByID`

## Middleware Global (handler.go)

```go
func (s *Handler) SetupGlobalMiddleware(handler http.Handler) http.Handler {
    handler = cors.New(cors.Options{
        AllowedHeaders:   []string{"*"},
        AllowedOrigins:   []string{"*"},
        AllowCredentials: true,
    }).Handler(handler)
    // wrapping de logging aqui
    return handler
}
```

## Autenticação

- JWT com RSA — validado pelo go-utils (`auth.AuthUser`)
- `principal *auth.AuthUser` chega como segundo argumento nas operations autenticadas
- Config: `APP_PUBLIC_KEY`, `APP_SECRET_KEY`, `APP_ISSUER_NAME`

## Tracking Info

Toda escrita propaga automaticamente `created_by`/`updated_by` a partir do `AuthUser` via `BaseModel.TrackingInfo` (JSONB).
