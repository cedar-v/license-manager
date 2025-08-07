# Multi-stage build
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM golang:1.23-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN go build -o license-manager cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend-builder /app/backend/license-manager .
COPY --from=frontend-builder /app/frontend/dist ./static
COPY backend/configs/config.example.yaml ./config.yaml
COPY backend/configs/i18n/ ./configs/i18n/
EXPOSE 18888
CMD ["./license-manager"]