FROM node:25.7.0-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

FROM golang:1.26.0 AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server -trimpath -ldflags="-s -w" ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=backend-builder /app/bin/server .
COPY --from=frontend-builder /app/dist ./frontend/dist

EXPOSE 8080

CMD ["./server"]
