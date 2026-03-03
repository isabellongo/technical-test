# Validação do banco (Fase 1.4)

## Conectar ao Postgres

```bash
# Via Docker (dentro do container)
docker exec -it driva-postgres psql -U postgres -d driva

# Via cliente local (psql, DBeaver, etc.)
# Host: localhost | Port: 5432 | User: postgres | Password: postgres | DB: driva
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
