# Banco de dados — inicialização e validação

Este diretório contém o `init.sql` usado para criar as tabelas do data warehouse (Bronze, Gold) e, opcionalmente, a tabela seed que popula a API simulada.

## Requisitos

- Docker & Docker Compose (o `docker-compose.yml` da raiz já expõe o serviço `postgres`).

## Como subir o Postgres (com compose)

```powershell
docker-compose up -d postgres
```

O serviço do compose normalmente é nomeado `postgres`. O container pode aparecer com um nome como `driva-postgres` dependendo do `docker-compose.yml`.

### Inicialização automática (init.sql)

O arquivo `db/init.sql` é montado no container em `/docker-entrypoint-initdb.d/01-init.sql` e será executado automaticamente apenas na primeira inicialização do volume do Postgres. Para reaplicar o script desde zero, remova o volume do Postgres e suba novamente:

```powershell
docker-compose down -v
docker-compose up -d postgres
```

Observação: remover o volume apaga todos os dados do banco — use com cautela.

## Como conectar ao Postgres

Via Docker (dentro do container):

```powershell
docker exec -it <nome_do_container_postgres> psql -U postgres -d driva
```

Via cliente local (psql, DBeaver, DataGrip):

- Host: localhost
- Porta: 5432
- User: postgres
- Password: <POSTGRES_PASSWORD> (ver `.env`)
- Database: driva

## Estrutura esperada das tabelas (resumo)

As tabelas principais criadas por `init.sql` (nomes sugestivos):

- `api_enrichments_seed` — dados seed que simulam milhares de enrichments (usada pela API fonte).
- `dw_bronze_enrichments` — camada Bronze: dados brutos ingeridos (camada de captura fiel).
- `dw_gold_enrichments` — camada Gold: dados transformados, com colunas em português e campos calculados.
- `dw_pipeline_state` (opcional) — controle do pipeline (última página processada, watermark, logs simples).

Consulte o `db/init.sql` para a definição completa de colunas e índices.

## Queries úteis para validação

```sql
-- listar tabelas
\dt

-- contar registros no seed (deve ser grande o suficiente para testar paginação)
SELECT COUNT(*) FROM api_enrichments_seed;

-- amostra do seed
SELECT id, workspace_name, total_contacts, contact_type, status
FROM api_enrichments_seed LIMIT 5;

-- validar ingestão Bronze
SELECT COUNT(*) FROM dw_bronze_enrichments;

-- validar processamento Gold
SELECT COUNT(*) FROM dw_gold_enrichments;
SELECT status_processamento, COUNT(*) FROM dw_gold_enrichments GROUP BY status_processamento;
```
Local do script: `db/init.sql`
