# React + TypeScript + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) uses [Babel](https://babeljs.io/) (or [oxc](https://oxc.rs) when used in [rolldown-vite](https://vite.dev/guide/rolldown)) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## React Compiler

The React Compiler is not enabled on this template because of its impact on dev & build performances. To add it, see [this documentation](https://react.dev/learn/react-compiler/installation).

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type-aware lint rules:
# Frontend — Dashboard (React + Vite)

Este diretório contém o dashboard de demonstração (React + TypeScript + Vite) que consome os endpoints de analytics da API.

## Requisitos técnicos

- Node.js 18+ e npm (ou yarn/pnpm)
- (Opcional) Docker, caso queira subir via container

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

## Rodar via Docker (opcional)

Se preferir executar o frontend em um container:

```powershell
docker build -f frontend/Dockerfile -t driva-frontend:local frontend
docker run --rm -p 5173:5173 -e VITE_API_URL=http://host.docker.internal:3000 driva-frontend:local
```

Observação: em Windows, `host.docker.internal` permite que o container acesse serviços expostos no host (ex.: API em `localhost`). Quando em Compose, prefira usar a rede do Compose e o serviço `api`.

## Testes / Verificação do funcionamento

- A aplicação não contém testes automatizados neste repositório (sugestão: adicionar testes com React Testing Library).
- Testes manuais de verificação:
  - Abra `http://localhost:5173`
  - Verifique a página principal com KPIs (total de enriquecimentos, % sucesso, tempo médio)
  - Navegue para a tabela/lista de enrichments e confirme paginação/filtros

Exemplo: se a API estiver rodando localmente e populada via n8n, as chamadas no frontend devem retornar dados visíveis.

## Dicas de desenvolvimento

- Ajuste a URL da API em `frontend/src/services/api.ts` através da variável `import.meta.env.VITE_API_URL`.
- Para hot-reload garantir que o Vite foi iniciado com a variável correta.
- Para produção, configure `VITE_API_URL` no processo de build.

## Melhorias sugeridas

- Adicionar testes unitários e de integração para os componentes-chave (KPIs, tabelas).
- Adicionar tratamento de erros na UI (mensagens claras ao usuário quando API retorna 429/erro).
- Adicionar storybook para documentar componentes UI.

---
Arquivo principal do frontend: `frontend/src/services/api.ts` (configurar URL da API).
        tsconfigRootDir: import.meta.dirname,
