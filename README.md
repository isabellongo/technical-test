# technical-test

Repositório de exemplo para o pipeline Driva (API Go, Postgres, n8n e frontend React/Vite).

## Visão geral

Pastas principais:
- `api/` - API em Go (Gin) que fornece endpoints de enrichments e analytics.
- `db/` - Scripts e instruções de inicialização do banco de dados.
- `n8n/` - Workflows n8n (JSON) para ingestão e processamento.
- `frontend/` - Aplicação React + Vite (UI de demonstração).

Para detalhes por componente, veja os READMEs nas respectivas pastas (`api/README.md`, `db/README.md`, `n8n/README.md`, `frontend/README.md`).

## Pré-requisitos (nível global)

- Docker & Docker Compose (Desktop ou Engine com Compose v2).
- Go 1.21+ (para rodar a API sem Docker, opcional se usar o container).
- Node.js 18+ e npm (para rodar o frontend localmente sem Docker).
- PowerShell (Windows) — exemplos de comandos neste README usam PowerShell.

## Variáveis de ambiente (.env)

Crie um arquivo `.env` na raiz do projeto com as variáveis mínimas abaixo (exemplo):

```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=driva
API_PORT=3000
API_KEY=driva_test_key_abc123xyz789
N8N_PORT=5678
```

O `docker-compose.yml` do projeto referencia essas variáveis. Se você usar os containers, o Postgres será exposto em `localhost:5432` e a API em `localhost:3000` por padrão.

## Como subir o projeto (modo rápido com Docker)

1. Garanta que o `.env` exista na raiz (ver seção acima).
2. Suba os serviços essenciais:

```powershell
docker-compose up -d postgres api n8n
```

3. (Opcional) Suba também o frontend em desenvolvimento numa janela separada se preferir:

```powershell
cd frontend
npm install
npm run dev
```

4. Verifique logs / healthchecks:

```powershell
docker-compose ps
docker-compose logs -f api
```

## Rodar serviços individualmente (sem Docker)

- API (Go):

```powershell
cd api
go mod download
go run ./cmd/main.go
```

- Frontend (Vite):

```powershell
cd frontend
npm install
npm run dev
```

## Comandos úteis

- Parar todos os containers:

```powershell
docker-compose down
```

- Recriar Postgres (quando quiser reaplicar `db/init.sql`):

```powershell
docker-compose down -v
docker-compose up -d postgres
```

Depois de subir o Postgres, o arquivo `db/init.sql` é montado em `/docker-entrypoint-initdb.d/01-init.sql` e será executado automaticamente se o volume do Postgres estiver vazio (primeira inicialização).

---

Verifique os READMEs locais para comandos e detalhes específicos de cada componente.