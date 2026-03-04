# API — Enrichments & Analytics

API escrita em Go (Gin) que fornece:

- Um endpoint que simula a fonte de enrichments (`/people/v1/enrichments`) com paginação, API Key e comportamento de rate-limit simulado.
- Endpoints de analytics que leem a camada Gold (`/analytics/*`) para consumo do dashboard.

Este README padroniza as instruções de execução, testes e validação da API.

## Requisitos técnicos

- Go 1.21+ (opcional — só necessário se rodar localmente sem Docker)
- Docker & Docker Compose (para rodar Postgres, n8n e a API em containers)
- Postgres 14+ (o compose já cuida disso)

## Variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto (ou exporte variáveis) com pelo menos:

- POSTGRES_USER (ex: postgres)
- POSTGRES_PASSWORD (ex: postgres)
- POSTGRES_DB (ex: driva)
- API_PORT (ex: 3000)
- API_KEY (ex: driva_test_key_abc123xyz789)

Dentro do Docker Compose, o host do Postgres é `postgres`.

## Como rodar (modo rápido com Docker)

1. Garanta que o `.env` exista na raiz.
2. Suba os serviços necessários:

```powershell
docker-compose up -d postgres api
```

3. Verifique que a API está respondendo:

```powershell
(Invoke-WebRequest -Uri "http://localhost:3000/health").Content
```

Resposta esperada: {"status":"ok"}

## Como rodar localmente (sem Docker)

1. Suba um Postgres local ou via Docker Compose:

```powershell
docker-compose up -d postgres
```

2. Na raiz, crie `.env` com as variáveis mínimas.

3. Instale dependências e rode:

```powershell
cd api
go mod download
go run ./cmd/main.go
```

## Endpoints principais

Autenticação: `Authorization: Bearer {API_KEY}`

- GET /people/v1/enrichments
   - Query: `page` (default 1), `limit` (default 50, max 100)
   - Resposta: objeto com `meta` (current_page, items_per_page, total_items, total_pages) e `data` (array de enrichments)
   - Observação: endpoint pode retornar 429 em simulações; n8n deve implementar retry/backoff.

- GET /analytics/overview
   - Retorna KPIs (total de jobs, % sucesso, tempo médio, etc.) consultando a camada Gold.

- GET /analytics/enrichments
   - Listagem paginada/filtrável da Gold (filtros exemplares: `id_workspace`, `status_processamento`, período, `categoria_tamanho_job`).

## Testes manuais (PowerShell)

Substitua a variável `API_KEY` conforme seu `.env` se necessário.

Health check:

```powershell
(Invoke-WebRequest -Uri "http://localhost:3000/health").Content
```

Fonte — Enrichments (exemplo):

```powershell
(Invoke-WebRequest -Uri "http://localhost:3000/people/v1/enrichments?page=1&limit=10" -Headers @{Authorization="Bearer driva_test_key_abc123xyz789"}).Content
```

Se receber `429 Too Many Requests`, reenvie após um pequeno intervalo; o comportamento é intencional para testar retry.

Analytics — Overview:

```powershell
(Invoke-WebRequest -Uri "http://localhost:3000/analytics/overview" -Headers @{Authorization="Bearer driva_test_key_abc123xyz789"}).Content
```

Analytics — Enrichments (com filtros):

```powershell
(Invoke-WebRequest -Uri "http://localhost:3000/analytics/enrichments?page=1&limit=10&status_processamento=CONCLUIDO" -Headers @{Authorization="Bearer driva_test_key_abc123xyz789"}).Content
```

### Alternativa com curl.exe (Windows)

```powershell
curl.exe -H "Authorization: Bearer driva_test_key_abc123xyz789" "http://localhost:3000/people/v1/enrichments?page=1&limit=10"
```

## Observabilidade e logs

- A API expõe logs básicos (startup, conexões DB, erros). Em ambiente containerizado, verifique com `docker-compose logs -f api`.
- Para produção, agregue um logger estruturado e métricas (Prometheus/Grafana) — sugestão na seção de melhorias.

## Seed e banco

- Há uma tabela opcional `api_enrichments_seed` (preenchida via `db/init.sql`) usada para simular milhares de registros. A API pagina esses dados via SQL.

## Boas práticas e melhorias sugeridas

- Adicionar testes automatizados (unit e integração) para endpoints e transformações do Gold.
- Expor métricas (Prometheus) e healthchecks mais completas (readiness/liveness).
- Implementar contratos OpenAPI/Swagger.

## Local dos arquivos importantes

- Código: `api/`
- Dockerfile: `api/Dockerfile`
- Inicialização DB: `db/init.sql`

---
Versão deste README: padronizado para instruções de execução e testes manuais.
