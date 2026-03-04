# Frontend — Dashboard (React + Vite)

Este diretório contém o dashboard de demonstração (React + TypeScript + Vite) que consome os endpoints de analytics da API.

## Requisitos técnicos

- Node.js 18+ e npm (ou yarn/pnpm)
- Docker, caso queira subir via container

## Variáveis / Configuração

O frontend espera uma variável para apontar a API. No Vite a convenção comum é `VITE_API_URL`.

- Exemplo `.env` (na pasta `frontend` ou na raiz dependendo do setup):

```
VITE_API_URL=http://localhost:3000
```

No repositório, a configuração da chamada à API está em `frontend/src/services/api.ts` — ajuste `VITE_API_URL` se necessário.

## Como rodar em desenvolvimento

```powershell
cd frontend
npm install
npm run dev
```

O servidor de desenvolvimento roda em `http://localhost:5173` por padrão.

## Build e preview

```powershell
npm run build
npm run preview
```

## Rodar via Docker

Se preferir executar o frontend em um container:

```powershell
docker build -f frontend/Dockerfile -t driva-frontend:local frontend
docker run --rm -p 5173:5173 -e VITE_API_URL=http://host.docker.internal:3000 driva-frontend:local
```

Observação: em Windows, `host.docker.internal` permite que o container acesse serviços expostos no host (ex.: API em `localhost`). Quando em Compose, prefira usar a rede do Compose e o serviço `api`.

## Testes / Verificação do funcionamento

- Testes manuais de verificação:
  - Abra `http://localhost:5173`
  - Verifique a página principal com KPIs (total de enriquecimentos, % sucesso, tempo médio)
  - Navegue para a tabela/lista de enrichments e confirme paginação/filtros

Exemplo: se a API estiver rodando localmente e populada via n8n, as chamadas no frontend devem retornar dados visíveis.

---
Arquivo principal do frontend: `frontend/src/services/api.ts` (configurar URL da API).
        tsconfigRootDir: import.meta.dirname,
