FROM golang:alpine AS builder

ARG STORAGE
ENV STORAGE=${STORAGE}

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# main app
RUN go build -o /app/bin/server ./cmd/url-shortening/main.go

# migrations (условно)
RUN if [ "$STORAGE" = "postgres" ]; then \
    go build -o /app/bin/migrate ./cmd/migrations/main.go; \
    fi

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin/server /app/bin/server
COPY --from=builder /app/bin/migrate /app/bin/migrate  

COPY . .

CMD ["/bin/sh", "-c", "if [ \"$STORAGE\" = \"postgres\" ]; then /app/bin/migrate --migrate=up && /app/bin/server; else /app/bin/server; fi"]
