# n8n - Workflows Driva (Fase 3)

Orquestração do pipeline de ingestão (API → Bronze → Gold).

---

## Pré-requisitos

### 1. Subir os serviços

```powershell
docker-compose up -d postgres api n8n
```

> Se o n8n já rodou antes com SQLite (padrão) e você mudou para Postgres, pode ser necessário:  
> `docker-compose down` e `docker-compose up -d` para o n8n criar as tabelas no Postgres.

### 2. Acessar o n8n

- **URL:** http://localhost:5678
- **Login:** admin / admin (conforme docker-compose)

### 3. Criar credencial Postgres

1. No n8n: **Settings** (engrenagem) → **Credentials** → **Add Credential**
2. Busque **Postgres**
3. Preencha:
   - **Host:** `postgres` (nome do serviço no Docker)
   - **Database:** `driva` (ou o valor de `POSTGRES_DB` no .env)
   - **User:** `postgres` (ou `POSTGRES_USER`)
   - **Password:** senha do seu `.env` (`POSTGRES_PASSWORD`)
   - **Port:** `5432`
4. **Nome sugerido:** `Postgres Driva` (para bater com os workflows)
5. Salve

---

## Importar workflows

1. Menu (três pontinhos) → **Import from File**
2. Selecione os arquivos na pasta `n8n/`, **nessa ordem**:
   - `01-ingestao-api-bronze.json`
   - `02-processamento-bronze-gold.json`
   - `03-orquestrador.json`
3. Em cada workflow, nos nós **Postgres**:
   - Clique no nó → **Credential to connect with** → selecione `Postgres Driva`
4. No **Orquestrador**, nos nós **Execute Workflow**:
   - Se aparecer "Workflow not found", selecione manualmente **Driva - Ingestão API → Bronze** e **Driva - Processamento Bronze → Gold**

---

## Testes da Fase 3

### 1. Testar Ingestão manualmente

1. Abra o workflow **Driva - Ingestão API → Bronze**
2. Clique em **Execute Workflow** (ou no botão ▶ do node Manual Trigger)
3. Verifique no Postgres:
   ```sql
   SELECT COUNT(*) FROM dw_bronze_enrichments;
   ```
   Esperado: ~5000 registros

### 2. Testar Processamento manualmente

1. Abra **Driva - Processamento Bronze → Gold**
2. **Execute Workflow**
3. Verifique:
   ```sql
   SELECT COUNT(*) FROM dw_gold_enrichments;
   SELECT status_processamento, COUNT(*) FROM dw_gold_enrichments GROUP BY status_processamento;
   ```

### 3. Testar Orquestrador (manual)

1. Abra **Driva - Orquestrador**
2. **Execute Workflow** (simula uma execução; o trigger por tempo roda a cada 5 min)
3. Ou **ative** o workflow (toggle) para rodar automaticamente a cada 5 min

### 4. Validar API analytics após rodar os workflows

```powershell
(Invoke-WebRequest -Uri "http://localhost:3000/analytics/overview" -Headers @{Authorization="Bearer driva_test_key_abc123xyz789"}).Content
```

---

## Ordem dos workflows

1. **Ingestão** — busca dados na API, grava na Bronze  
2. **Processamento** — Bronze → Gold (transformações)  
3. **Orquestrador** — agenda Ingestão + Processamento a cada 5 min  

---

## Decisões técnicas

| Decisão | Escolha | Motivo |
|---------|---------|--------|
| n8n + Postgres | `DB_TYPE=postgresdb` | Workflows e execuções persistidas no banco |
| Paginação API | Code node com `fetch` | Um único node, retry em 429 e loop de páginas |
| Retry 429 | Loop `while` dentro do Code | Aguarda 2s e tenta de novo na mesma página |
| Bronze→Gold | Code node com regras em JS | Mapeamento PT, categorias, flags e duração |
| Orquestrador | Schedule 5 min + Execute Workflow | Ingestão e Processamento em sequência |
| Credencial Postgres | Nome `Postgres Driva` | Padrão usado nos workflows; host `postgres` na rede Docker |
| API URL no Code | `http://api:3000` | Nome do serviço na rede Docker Compose |

### Regras Bronze → Gold

- **status_processamento:** PROCESSING→EM_PROCESSAMENTO, COMPLETED→CONCLUIDO, FAILED→FALHOU, CANCELED→CANCELADO
- **tipo_contato:** PERSON→PESSOA, COMPANY→EMPRESA
- **categoria_tamanho_job:** 0–100 PEQUENO, 101–500 MEDIO, 501–1500 GRANDE, >1500 MUITO_GRANDE
- **processamento_sucesso:** `true` só se status = COMPLETED
- **necessita_reprocessamento:** `true` se status = FAILED ou PROCESSING
