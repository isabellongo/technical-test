-- =============================================================================
-- init.sql - Inicialização do Data Warehouse Driva
-- =============================================================================
-- DECISÃO: Este arquivo é executado automaticamente pelo Postgres na primeira
-- subida do container (docker-entrypoint-initdb.d). Em produção, normalmente
-- usa-se migrations (Flyway, Liquibase, golang-migrate) para controle de versão.
-- Aqui usamos init.sql para simplicidade do ambiente de teste.
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 1. api_enrichments_seed
-- -----------------------------------------------------------------------------
-- AMBIENTE DE TESTE: Tabela que simula a API externa de enriquecimentos.
-- Em produção, esses dados viriam de uma API real; aqui servimos paginado
-- para o n8n consumir e testar o pipeline ponta-a-ponta.
--
-- DECISÃO: Estrutura espelha exatamente o JSON da API (GET /people/v1/enrichments)
-- para facilitar a implementação da API (SELECT + paginação simples).
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS api_enrichments_seed (
    id                  UUID PRIMARY KEY,
    id_workspace        UUID NOT NULL,
    workspace_name      VARCHAR(255) NOT NULL,
    total_contacts      INTEGER NOT NULL CHECK (total_contacts >= 0),
    contact_type        VARCHAR(50) NOT NULL,  -- API: "PERSON" | "COMPANY"
    status              VARCHAR(50) NOT NULL,  -- API: "PROCESSING" | "COMPLETED" | "FAILED" | "CANCELED"
    created_at          TIMESTAMPTZ NOT NULL,
    updated_at          TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_seed_workspace ON api_enrichments_seed(id_workspace);
CREATE INDEX IF NOT EXISTS idx_seed_status ON api_enrichments_seed(status);
CREATE INDEX IF NOT EXISTS idx_seed_created ON api_enrichments_seed(created_at);

-- Necessário para gen_random_uuid() usado no seed
CREATE EXTENSION IF NOT EXISTS pgcrypto;


-- -----------------------------------------------------------------------------
-- 2. dw_bronze_enrichments (Camada Bronze)
-- -----------------------------------------------------------------------------
-- AMBIENTE NORMAL: A Bronze armazena os dados brutos tal qual recebidos da fonte.
-- Princípio "captura fiel" — não alteramos estrutura ou tipos.
--
-- DECISÃO: Usamos tipos nativos (VARCHAR, INTEGER, TIMESTAMPTZ) em vez de JSONB
-- porque a API retorna estrutura fixa. JSONB seria útil se a API fosse schema-less
-- ou variasse muito; para este desafio, tipos explícitos facilitam upsert e índices.
--
-- Campos de controle DW: dw_ingested_at (primeira ingestão), dw_updated_at
-- (última atualização) — exigidos pelo desafio para rastreabilidade.
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS dw_bronze_enrichments (
    id                  UUID PRIMARY KEY,
    id_workspace        UUID NOT NULL,
    workspace_name      VARCHAR(255) NOT NULL,
    total_contacts      INTEGER NOT NULL,
    contact_type        VARCHAR(50) NOT NULL,
    status              VARCHAR(50) NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL,
    updated_at          TIMESTAMPTZ NOT NULL,
    dw_ingested_at      TIMESTAMPTZ NOT NULL,
    dw_updated_at       TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_bronze_id ON dw_bronze_enrichments(id);
CREATE INDEX IF NOT EXISTS idx_bronze_workspace ON dw_bronze_enrichments(id_workspace);
CREATE INDEX IF NOT EXISTS idx_bronze_updated ON dw_bronze_enrichments(dw_updated_at);


-- -----------------------------------------------------------------------------
-- 3. dw_gold_enrichments (Camada Gold)
-- -----------------------------------------------------------------------------
-- AMBIENTE NORMAL: Gold = dados processados, prontos para analytics e dashboard.
-- Colunas em português e campos calculados conforme regras de negócio.
--
-- DECISÕES:
-- - Nomes em PT: requisito do desafio; facilita leitura pelo time de Visibilidade.
-- - Campos calculados pré-computados aqui (não na view) para performance do
--   dashboard e queries analíticas.
-- - id_enriquecimento = id original, usado como chave de upsert.
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS dw_gold_enrichments (
    id_enriquecimento           UUID PRIMARY KEY,
    id_workspace                UUID NOT NULL,
    nome_workspace              VARCHAR(255) NOT NULL,
    total_contatos              INTEGER NOT NULL,
    tipo_contato                VARCHAR(50) NOT NULL,   -- PESSOA | EMPRESA
    status_processamento        VARCHAR(50) NOT NULL,   -- EM_PROCESSAMENTO | CONCLUIDO | FALHOU | CANCELADO
    data_criacao                TIMESTAMPTZ NOT NULL,
    data_atualizacao            TIMESTAMPTZ NOT NULL,
    duracao_processamento_minutos FLOAT,
    tempo_por_contato_minutos   FLOAT,
    processamento_sucesso       BOOLEAN NOT NULL,
    categoria_tamanho_job       VARCHAR(20) NOT NULL,   -- PEQUENO | MEDIO | GRANDE | MUITO_GRANDE
    necessita_reprocessamento   BOOLEAN NOT NULL,
    data_atualizacao_dw         TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_gold_id ON dw_gold_enrichments(id_enriquecimento);
CREATE INDEX IF NOT EXISTS idx_gold_workspace ON dw_gold_enrichments(id_workspace);
CREATE INDEX IF NOT EXISTS idx_gold_status ON dw_gold_enrichments(status_processamento);
CREATE INDEX IF NOT EXISTS idx_gold_data_criacao ON dw_gold_enrichments(data_criacao);
CREATE INDEX IF NOT EXISTS idx_gold_categoria ON dw_gold_enrichments(categoria_tamanho_job);


-- -----------------------------------------------------------------------------
-- 4. dw_pipeline_state (Controle / Watermark)
-- -----------------------------------------------------------------------------
-- AMBIENTE NORMAL: Usado para ingestão incremental ou estado do pipeline.
-- AMBIENTE DE TESTE: MVP usa full refresh (busca todas as páginas). Esta tabela
-- permite evoluir para watermark (ex: última página/updated_at processada).
--
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS dw_pipeline_state (
    id                  SERIAL PRIMARY KEY,
    pipeline_name       VARCHAR(100) NOT NULL UNIQUE,
    last_page_processed INTEGER,
    last_run_at         TIMESTAMPTZ,
    metadata            JSONB,
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO dw_pipeline_state (pipeline_name, last_page_processed, last_run_at)
VALUES ('ingest_enrichments', NULL, NULL)
ON CONFLICT (pipeline_name) DO NOTHING;


-- =============================================================================
-- SEED: ~5000 registros para simular a API
-- =============================================================================
-- AMBIENTE DE TESTE: Dados sintéticos para validar paginação, transformações
-- e pipeline. Em produção, a fonte seria a API real de enriquecimentos.
--
-- Distribuição variada:
-- - status: ~60% COMPLETED, ~15% PROCESSING, ~15% FAILED, ~10% CANCELED
-- - contact_type: ~50% COMPANY, ~50% PERSON
-- - total_contacts: 10 a 2000 (permite todas as categorias: PEQUENO, MEDIO, GRANDE, MUITO_GRANDE)
-- - workspaces: ~20 workspaces diferentes
-- =============================================================================

INSERT INTO api_enrichments_seed (id, id_workspace, workspace_name, total_contacts, contact_type, status, created_at, updated_at)
SELECT
    gen_random_uuid(),
    (array['e6bb64bf-46e4-410d-8406-c61e267ea607','a1b2c3d4-e5f6-7890-abcd-ef1234567890','b2c3d4e5-f6a7-8901-bcde-f12345678901','c3d4e5f6-a7b8-9012-cdef-123456789012','d4e5f6a7-b8c9-0123-defa-234567890123','e5f6a7b8-c9d0-1234-efab-345678901234','f6a7b8c9-d0e1-2345-fabc-456789012345','a7b8c9d0-e1f2-3456-abcd-567890123456','b8c9d0e1-f2a3-4567-bcde-678901234567','c9d0e1f2-a3b4-5678-cdef-789012345678','d0e1f2a3-b4c5-6789-defa-890123456789','e1f2a3b4-c5d6-7890-efab-901234567890','f2a3b4c5-d6e7-8901-fabc-a12345678901','a3b4c5d6-e7f8-9012-abcd-b23456789012','b4c5d6e7-f8a9-0123-bcde-c34567890123','c5d6e7f8-a9b0-1234-cdef-d45678901234','d6e7f8a9-b0c1-2345-defa-e56789012345','e7f8a9b0-c1d2-3456-efab-f67890123456','f8a9b0c1-d2e3-4567-fabc-a78901234567','a9b0c1d2-e3f4-5678-abcd-b89012345678']::uuid[])[1 + (i % 20)],
    'Workspace ' || (1 + (i % 20)),
    10 + (random() * 1990)::integer,
    (array['COMPANY','PERSON'])[1 + (i % 2)],
    (array['COMPLETED','COMPLETED','COMPLETED','PROCESSING','FAILED','CANCELED'])[1 + (i % 6)],
    NOW() - (random() * interval '90 days'),
    NOW() - (random() * interval '89 days')
FROM generate_series(1, 5000) AS i;

-- Nota: O init.sql roda apenas na primeira subida do container (volume vazio).
-- Para resetar: docker-compose down -v && docker-compose up -d postgres
