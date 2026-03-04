# n8n — Workflows do pipeline (Ingestão, Processamento e Orquestração)

Este diretório contém os exports dos workflows do n8n usados para:

- Ingestão: buscar a API de enrichments e persistir na camada Bronze.
- Processamento: transformar registros da Bronze em Gold (camada analítica).
- Orquestrador: agendar (5 minutos) e encadear ingestão → processamento.

## Requisitos

- Docker & Docker Compose (o `docker-compose.yml` já referencia o serviço `n8n`).
- Acesso ao serviço: http://localhost:5678 quando em Docker Compose.

## Como subir o n8n (com Docker Compose)

```powershell
docker-compose up -d postgres api n8n
```

Se já houver dados antigos em SQLite e você mudou para Postgres, faça `docker-compose down` e depois `docker-compose up -d` para criar as tabelas no Postgres.

## Acessando o n8n

- URL: http://localhost:5678
- Credenciais padrão (quando criado via compose): `admin` / `admin` — verifique `docker-compose.yml` para valores.

## Criar credencial Postgres no n8n

1. Settings → Credentials → Add Credential → Postgres
2. Preencha com:
   - Host: `postgres`
   - Database: `driva`
   - User: `postgres` (ou `POSTGRES_USER`)
   - Password: `POSTGRES_PASSWORD`
   - Port: `5432`
3. Nome recomendado: `Postgres Driva`

## Importar workflows

1. Menu (três pontinhos) → Import from File
2. Importe, na ordem abaixo, os arquivos contidos nesta pasta:
   - `01-ingestao-api-bronze.json`
   - `02-processamento-bronze-gold.json`
   - `03-orquestrador.json`
3. Para cada nó Postgres, selecione a credential criada (`Postgres Driva`).
4. No Orquestrador, valide os nós Execute Workflow apontando para os workflows importados.

## Testes / Validação (passos rápidos)

1. Testar ingestão manual:
   - Abra `Driva - Ingestão API → Bronze` e clique em Execute (Manual Trigger).
   - Verifique no Postgres:

   ```sql
   SELECT COUNT(*) FROM dw_bronze_enrichments;
   ```

   Esperado: registros importados (seed padrão ~5000).

2. Testar processamento manual:
   - Abra `Driva - Processamento Bronze → Gold` e Execute.
   - Verifique:

   ```sql
   SELECT COUNT(*) FROM dw_gold_enrichments;
   SELECT status_processamento, COUNT(*) FROM dw_gold_enrichments GROUP BY status_processamento;
   ```

3. Testar orquestrador:
   - Abra `Driva - Orquestrador` e execute manualmente para simular o cron.
   - Ou ative (toggle) para que rode a cada 5 minutos.

4. Validar API Analytics:

```powershell
(Invoke-WebRequest -Uri "http://localhost:3000/analytics/overview" -Headers @{Authorization="Bearer driva_test_key_abc123xyz789"}).Content
```

## Observabilidade e boas práticas nos workflows

- Implementar logging (nodes de Set/Function que gravam contagem/erros em uma tabela `dw_pipeline_state`).
- Retry em 429 (HTTP Request): backoff exponencial, 5 tentativas é um bom ponto inicial.
- Usar Execuções chamáveis (Execute Workflow) para manter ingestão e processamento modulares.

## Estrutura e decisões técnicas (resumo)

- DB Type: `postgresdb` (persistência de workflows/executions no Postgres).
- Paginação da API: loop via code node (`this.helpers.httpRequest`) com tratamento de 429.
- Bronze→Gold: Code node com regras de tradução e cálculos (duracao, tempo_por_contato, flags).
- Orquestrador: agendador 5 minutos → chama ingestão → quando concluído chama processamento.

## Arquivos nesta pasta

- `01-ingestao-api-bronze.json` — workflow de ingestão (importar primeiro)
- `02-processamento-bronze-gold.json` — workflow de transformação
- `03-orquestrador.json` — scheduler e encadeamento

---
Notas: se alterar URLs (ex.: `http://api:3000`), atualize os code nodes que fazem as requisições.
