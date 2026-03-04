# Driva Data Pipeline & Analytics Challenge

Este repositório contém a solução completa para o desafio técnico da **Driva**. A aplicação consiste em um ecossistema conteinerizado que realiza a ingestão, processamento e visualização de dados de enriquecimento B2B, seguindo a arquitetura de **Medallion Data Warehouse (Bronze & Gold)**.

## Arquitetura do Sistema

A solução foi desenhada para ser resiliente e escalável, utilizando o **n8n** como orquestrador de ETL, **Go** para a camada de serviços e **React** para a interface de monitoramento.

```mermaid
flowchart LR

subgraph Fonte [Fonte Simulada]
API_Fonte[API /people/v1/enrichments]
end

subgraph Ingestao [Ingestão]
    n8n_Ingest[n8n Workflow Ingestão]
end

subgraph DW [Data Warehouse]
    Bronze[(Bronze - dados brutos)]
    Gold[(Gold - dados processados)]
end

subgraph Processamento [Processamento]
    n8n_Proc[n8n Workflow Bronze→Gold]
end

subgraph Orquestrador [Orquestrador]
    n8n_Sched[n8n Scheduler 5min]
end

subgraph Consumo [Consumo]
    API_Analytics[API /analytics/*]
    Dashboard[Dashboard React]
end

n8n_Sched --> n8n_Ingest
n8n_Ingest --> API_Fonte
API_Fonte -->|"HTTP GET paginado"| n8n_Ingest
n8n_Ingest -->|"Insert/Upsert"| Bronze

n8n_Sched --> n8n_Proc
n8n_Proc -->|"SELECT"| Bronze
n8n_Proc -->|"Upsert"| Gold

Dashboard -->|"GET /analytics/*"| API_Analytics
API_Analytics -->|"SELECT"| Gold
```

## Mapa de Implementação
O projeto foi executado seguindo um cronograma de marcos técnicos para garantir a integração contínua de cada camada:


```mermaid
flowchart TD
A[0. Prep: GitHub, pastas, .env] --> B[1. init.sql + docker-compose]
B --> C[2. API Go - endpoints fonte]
C --> D[2b. API - endpoints analytics]
D --> E[3. Workflows n8n]
E --> F[4. Dashboard React]
F --> G[5. README + Vídeo]
G --> H{MVP completo?}
H -->|Sim| I[Melhorias opcionais]
```
## Decisões de Engenharia

1. **Idempotência e Resiliência** (Cláusula UPSERT): A lógica de Upsert nas camadas de dados garante que não haverá duplicidade.
2. **Separação de preocupações:** A camada **Bronze** armazena o dado bruto (Raw), para que mudanças nas regras de negócio entre Bronze -> Gold não precisem buscar todos os dados na API novamente.
3. **Estratégia de Paginação e Gestão de Memória:** Para garantir que o projeto funcione com qualquer quantidade de entradas, foram aplicadas as seguintes soluções: 
- Prevenção de Buffer Overflow: Em vez de tentar carregar todos os dados em memória para uma única inserção, o pipeline processa "página por página". Isso mantém o consumo de RAM do container n8n constante e baixo, independente do tamanho do dataset.
- Iteração Sequencial vs. Rate Limit: A escolha pela paginação sequencial (em vez de disparos paralelos) foi uma decisão deliberada para respeitar o Rate Limiting da API e garantir que a ordem dos logs de ingestão (dw_ingested_at) seja consistente.
4. **Estratégia de Backoff (Tratamento de 429):** A API fonte simula limites de requisição.
- Uma lógica de **Wait/Retry** que atua como um _Exponential Backoff_ simplificado foi utilizada, evitando que o pipeline entre em um loop de erro infinito e "estresse" o servidor da API.
  
## Boas práticas aplicadas:
- **Separação Clara de Responsabilidades:** Arquitetura limpa com divisão entre lógica de domínio (internal), rotas e repositórios na API Go.

- **Conteinerização Full-Stack:** Todo o ambiente (Postgres, n8n, API e Frontend) sobe via Docker, garantindo paridade de ambiente e facilidade de deploy (Publishing).

- **Resiliência no Pipeline (n8n):** Tratamento de Rate Limit (429 Too Many Requests) com políticas de Retry e Backoff.

	- Paginação inteligente via loop dinâmico que detecta o fim dos registros na origem.

	- Mecanismo de controle para evitar disparos simultâneos e garantir a integridade da sequência de páginas.

- **Modelagem Medallion** (Bronze/Gold): 
	- **Bronze:** Captura fiel (Raw) com campos de controle dw_ingested_at e dw_updated_at.
	- **Gold:** Camada de negócio com transformações de nomes (PT-BR), tradução de status, cálculos de duração de processamento e segmentação de jobs por tamanho.

- **Segurança:** Implementação de Middleware de autenticação via API Key (Bearer Token) em todos os endpoints sensíveis.

- **UI/UX Identitária:** Dashboard desenvolvido com React + Vite e Tailwind CSS, aplicando a paleta de cores institucional da Driva.

- **Versionamento e Testes:** Uso de Git para versionamento e testes realizados para cada etapa do projeto antes da integração final (detalhados nos sub-readmes).

## Como Executar
Mais detalhes sobre os pré-requisitos individuais e os testes para verificação de cada etapa estão disponíveis nos README.md dentro das pastas /db, /api, /frontend, /n8n.

**Pré-requisitos:** 
- Docker & Docker Compose (v2+).
- Arquivo .env na raiz (conforme exemplo abaixo).

**Configuração Rápida**
1. **Variáveis de Ambiente:** Crie o arquivo .env na raiz do projeto:
```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=driva
API_PORT=3000
API_KEY=driva_test_key_abc123xyz789
N8N_PORT=5678
```

2. **Subir Infraestrutura:**
```
docker-compose up -d
```

3. **Configurar o n8n**
	1. Acesse http://localhost:5678.
	2. Importe os arquivos JSON localizados na pasta /n8n.
	3. Configure as credenciais do Postgres e a API Key (driva_test_key_abc123xyz789).
	4. Execute o Orquestrador para iniciar a primeira carga.

**Portas dos Serviços**
- Frontend (Docker): http://localhost:5173
- API (Go): http://localhost:3000
- n8n: http://localhost:5678
- Postgres: http://localhost:5432

**Testando os Endpoints**
1. Consumo da Fonte (paginado):
```
curl -H "Authorization: Bearer driva_test_key_abc123xyz789" \
"http://localhost:3000/people/v1/enrichments?page=1&limit=10"
```
2. Analytics Overview: 
```
curl -H "Authorization: Bearer driva_test_key_abc123xyz789" \
"http://localhost:3000/analytics/overview"
```
3. Listagem de Enriquecimentos:
```
curl -H "Authorization: Bearer driva_test_key_abc123xyz789" \
"http://localhost:3000/analytics/enrichments?limit=5"
```

## Estrutura do Repositório
Para detalhes específicos de cada componente, consulte os arquivos README internos:
- /api: API em Go (Gin), middlewares e endpoints de analytics.
- /frontend: Aplicação React + Vite, hooks e componentes de gráfico.
- /n8n: JSONs dos workflows e guia de importação.
- /db: Script init.sql e instruções de seed do banco de dados.
## Demonstração
Um vídeo explicando a arquitetura, o funcionamento do pipeline e a visualização dos dados no Dashboard pode ser encontrado no link abaixo:

👉 [LINK_PARA_O_VIDEO_AQUI]

*Desenvolvido como parte do desafio técnico para o time de Tech da Driva.*
