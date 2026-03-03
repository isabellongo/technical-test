# Validação do banco (Fase 1.4)

## Pré-requisitos

- Docker & Docker Compose (para subir o serviço Postgres via `docker-compose.yml`).

## Subir Postgres via Docker Compose

```powershell
docker-compose up -d postgres
```

O serviço Postgres no compose é nomeado `postgres` e o container terá o nome `driva-postgres`.

### Importante: inicialização do banco

O arquivo `db/init.sql` está montado no container em `/docker-entrypoint-initdb.d/01-init.sql` e será executado automaticamente na primeira inicialização do banco (quando o volume estiver vazio). Se precisar reaplicar o script desde o início, remova o volume do Postgres e suba novamente:

```powershell
docker-compose down -v
docker-compose up -d postgres
```

## Conectar ao Postgres

```powershell
# Via Docker (dentro do container)
docker exec -it driva-postgres psql -U postgres -d driva

# Via cliente local (psql, DBeaver, etc.)
# Host: localhost | Port: 5432 | User: postgres | Password: <POSTGRES_PASSWORD> | DB: driva
```

## Queries de validação

```sql
-- Verificar tabelas
\dt

-- Contar registros no seed
SELECT COUNT(*) FROM api_enrichments_seed;

-- Amostra do seed
SELECT id, workspace_name, total_contacts, contact_type, status
FROM api_enrichments_seed LIMIT 5;

-- Verificar Bronze e Gold (vazios até rodar os workflows)
SELECT COUNT(*) FROM dw_bronze_enrichments;
SELECT COUNT(*) FROM dw_gold_enrichments;
```
