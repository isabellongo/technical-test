# API Driva - Fase 2

API Go + Gin que expõe a fonte de enrichments (simulada) e endpoints de analytics.

---

## Pré-requisitos

### 1. Instalar Go 1.21+

- **Windows:** Baixe em [golang.org/dl](https://go.dev/dl/) e execute o instalador.
- **Verificar:** `go version`

### 2. Instalar dependências do projeto

Na pasta `api/`:

```powershell
cd api
go mod download
```

Ou, para baixar e atualizar `go.sum`:

```powershell
go mod tidy
```

---

## Variáveis de ambiente

Defina na raiz do projeto (arquivo `.env`) ou exporte no terminal:

| Variável        | Descrição                     | Exemplo                    |
|-----------------|-------------------------------|----------------------------|
| POSTGRES_HOST   | Host do Postgres              | `localhost` (local) / `postgres` (Docker) |
| POSTGRES_PORT   | Porta do Postgres             | `5432`                     |
| POSTGRES_USER   | Usuário                       | `postgres`                 |
| POSTGRES_PASSWORD | Senha                       | `changeme`                 |
| POSTGRES_DB     | Nome do banco                 | `driva`                    |
| API_PORT        | Porta da API                  | `3000`                     |
| API_KEY         | Chave para Authorization      | `driva_test_key_abc123xyz789` |

Para Docker, o compose injeta as variáveis automaticamente. `POSTGRES_HOST` é definido como `postgres` dentro da rede.

---

## Rodar localmente (sem Docker)

1. Subir o Postgres:
   ```powershell
   docker-compose up -d postgres
   ```

2. Criar `.env` na raiz (se ainda não tiver) e garantir `API_KEY`:
   ```
   API_KEY=driva_test_key_abc123xyz789
   ```

3. Rodar a API:
   ```powershell
   cd api
   go run ./cmd/main.go
   ```

---

## Rodar via Docker

```powershell
docker-compose up -d postgres api
```

---

## Testes da Fase 2

### Health check

```powershell
curl http://localhost:3000/health
```

Resposta esperada: `{"status":"ok"}`

### Fonte - Enrichments paginados

```powershell
curl -H "Authorization: Bearer driva_test_key_abc123xyz789" "http://localhost:3000/people/v1/enrichments?page=1&limit=10"
```

(Em ~5% das requisições pode retornar 429 — simulação para teste de retry no n8n.)

### Analytics - Overview (KPIs)

```powershell
curl -H "Authorization: Bearer driva_test_key_abc123xyz789" "http://localhost:3000/analytics/overview"
```

Nota: A Gold pode estar vazia até rodar os workflows n8n. Nesse caso, totais serão 0.

### Analytics - Enrichments paginados

```powershell
curl -H "Authorization: Bearer driva_test_key_abc123xyz789" "http://localhost:3000/analytics/enrichments?page=1&limit=10"
```

Com filtros (opcionais):
```powershell
curl -H "Authorization: Bearer driva_test_key_abc123xyz789" "http://localhost:3000/analytics/enrichments?page=1&limit=10&status_processamento=CONCLUIDO&categoria_tamanho_job=MEDIO"
```

---

## Teste sem chave (esperado: 401)

```powershell
curl http://localhost:3000/people/v1/enrichments
```

Resposta esperada: `{"error":"missing Authorization header"}`

---

## Decisões técnicas

| Decisão          | Escolha            | Motivo                                          |
|------------------|--------------------|--------------------------------------------------|
| Driver DB        | pgx/v5             | Driver oficial e mantido, suporte a UUID nativo  |
| Framework        | Gin                | Conforme plano; leve, middleware pronto          |
| Simulação 429    | 5% por requisição  | Endpoint `/people/v1/enrichments` apenas; para testar retry no n8n |
| API Key          | Bearer no header   | `Authorization: Bearer <key>`; valor do .env com fallback para `driva_test_key_abc123xyz789` |
| Paginação        | page, limit        | default 50, max 100                              |
| Dockerfile       | Multi-stage Alpine | Imagem final pequena                             |
