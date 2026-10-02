FROM golang:1.27.1-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/main.go


FROM node:24.21.0-alpine AS oas2pm
WORKDIR /tools/oas2pm
COPY tools/oas2pm/package.json tools/oas2pm/package-lock.json ./
RUN npm ci --omit=dev
COPY tools/oas2pm/convert.js ./


FROM node:24.21.0-alpine
WORKDIR /app

ENV TZ=Europe/Moscow
RUN apk add --no-cache ca-certificates tzdata curl

COPY --from=oas2pm /tools/oas2pm ./tools/oas2pm
COPY --from=builder /out/app ./app

USER node
CMD ["sleep", "infinity"]