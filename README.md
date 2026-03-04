# Driva Data Pipeline & Analytics Challenge

Este repositório contém a solução completa para o desafio técnico da **Driva**. A aplicação consiste em um ecossistema conteinerizado que realiza a ingestão, processamento e visualização de dados de enriquecimento B2B, seguindo a arquitetura de **Medallion Data Warehouse (Bronze & Gold)**.

## Arquitetura do Sistema
* **Go (Gin):** Escolhido pela alta performance em concorrência e baixo footprint de memória, ideal para serviços de alto volume de dados e APIs de baixa latência.
* **n8n:** Utilizado como motor de orquestração por permitir rápida iteração em fluxos de ETL complexos com gestão nativa de retries e agendamento.
* **PostgreSQL:** Banco de dados relacional robusto para garantir a integridade referencial e suporte nativo a operações de UPSERT e tipos JSONB.
* **React + Vite:** Stack moderna para o frontend que garante builds rápidos e uma interface reativa para visualização de métricas.

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
## Decisões de Engenharia e Boas Práticas

### Performance e Experiência do Desenvolvedor (DX)
* **Indexação Estratégica:** O banco de dados foi projetado com índices nas colunas de busca frequente (`id_workspace`, `status_processamento`, `data_atualizacao_dw`) na camada Gold, garantindo que o Dashboard escale mesmo com milhões de registros.
* **Determinismo de Testes:** O script `init.sql` foi construído para ser idempotente e inclui um seed determinístico, permitindo que qualquer desenvolvedor suba o ambiente e visualize os mesmos indicadores imediatamente.

### Confiabilidade de Dados (Idempotência e Medallion)
* **Cláusula UPSERT:** A lógica de ON CONFLICT nas camadas Bronze e Gold garante que o pipeline seja idempotente. Se o job falhar ou rodar em duplicidade, o banco apenas atualiza o estado atual em vez de gerar registros duplicados.
* **Arquitetura Medallion:** A camada Bronze captura o dado bruto (Raw) de forma fiel. Isso permite que novas regras de negócio na Gold sejam aplicadas retroativamente sem a necessidade de re-onerar a API fonte.

### Eficiência de Pipeline (Paginação e Backoff)
* **Gestão de Memória:** O pipeline processa os dados página por página em vez de carregar 5.000 registros na memória do n8n de uma vez, prevenindo Buffer Overflows.
* **Exponential Backoff:** Em caso de erro 429 (Too Many Requests), o workflow implementa uma lógica de espera inteligente, evitando o estresse desnecessário da infraestrutura.
* **Orquestração Atômica:** O Orquestrador garante a integridade sequencial: a camada Gold só inicia o processamento após a confirmação de sucesso da ingestão na Bronze.

### Qualidade de Software
* **Clean Architecture:** API Go estruturada com separação clara entre lógica de domínio (internal), repositórios e rotas.
* **Segurança:** Implementação de Middleware para validação de API Key (Bearer Token) em todos os endpoints sensíveis.
* **UI/UX Identitária:** Interface desenvolvida com Tailwind CSS utilizando a paleta de cores institucional da Driva para uma experiência de produto completa.

## Segurança
Durante a fase de desenvolvimento, testes e verificação inicial do sistema (como demonstrado nas configurações de ambiente deste repositório), algumas credenciais, senhas de teste e API Keys podem estar visíveis ou configuradas de forma estática.
Essas chaves e senhas expostas têm como **única** finalidade facilitar a **validação funcional** do sistema e garantir que os fluxos de integração (como n8n, Postgres e APIs externas) estejam operando conforme o esperado. Nenhum dos dados ou credenciais apresentados é real.

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

👉 https://youtu.be/GOc-447Um8Q

## Melhorias e Expansões Futuras
### Observabilidade e Saúde do Sistema
* **Métricas com Prometheus/Grafana:** Implementação de um exportador de métricas na API Go para monitorar latência de endpoints e taxa de erro 429 da fonte.
* **Logs Estruturados:** Substituição dos logs padrão por logs em JSON (ex: Zap ou Logrus) para facilitar a agregação em ferramentas como ELK Stack ou Datadog.

### Qualidade e Contratos
* **Testes Automatizados:** Implementação de testes de integração para as transformações Bronze→Gold e testes de componentes no Frontend utilizando React Testing Library.
* **Documentação de API:** Integração com Swagger/OpenAPI para gerar documentação interativa dos endpoints de Analytics.

### Frontend Pro
* **Tratamento de Erros na UI:** Adição de estados de "Loading" e "Toast Notifications" para informar o usuário caso a API fonte esteja em Rate Limit ou fora do ar.
* **Storybook:** Documentação visual dos componentes de KPI e Gráficos para garantir consistência visual.

### Watermark para Ingestão Incremental
Atualmente, o pipeline processa o dataset completo para garantir a integridade. Uma melhoria crítica (já com 50% da infraestrutura pronta no banco) é a implementação de Watermark. Isso permitirá que o n8n consulte apenas registros criados ou alterados desde a última execução bem-sucedida.

### Observabilidade e Logs
* Health Checks Dinâmicos: Implementação de um endpoint /health na API Go que verifica a conectividade ativa com o Postgres.
* Monitoramento de Workflows: Integração do n8n com ferramentas de log externas para alertar caso o workflow falhe consecutivamente.

###	Qualidade de Dados
* Camada de Validação: Adição de testes de integridade entre a Bronze e a Gold (ex: garantir que total_contacts nunca seja negativo) antes da transformação final.

*Desenvolvido por Isabel Cristina Kavalco Longo como parte do desafio técnico para o time de Tech da Driva.*
