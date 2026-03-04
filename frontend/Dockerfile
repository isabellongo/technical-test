FROM node:18-alpine AS builder
WORKDIR /app

# Instalar dependências
COPY package.json package-lock.json* ./
COPY tsconfig.app.json tsconfig.json ./
RUN npm ci --silent

# Copiar código e construir
COPY . .

# Build-time variables for Vite
ARG VITE_API_URL
ARG VITE_API_KEY
# For robust deployment with nginx proxy, default VITE_API_URL to '/api'
RUN echo "VITE_API_URL=${VITE_API_URL:-/api}" > .env
RUN echo "VITE_API_KEY=${VITE_API_KEY:-driva_test_key_abc123xyz789}" >> .env

RUN npm run build

FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD wget -qO- http://localhost/ >/dev/null || exit 1

CMD ["nginx", "-g", "daemon off;"]
