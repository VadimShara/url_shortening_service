FROM golang:alpine AS builder

ARG STORAGE
ENV STORAGE=${STORAGE}

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# main app
RUN go build -o /app/bin/server ./cmd/url-shortening/main.go

RUN if ["$STORAGE" = "postgres"]; then \
    # migrations \
    RUN go build -o /app/bin/migrate ./cmd/migrations/main.go; \
    fi

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin/server /app/bin/server
RUN if ["$STORAGE" = "postgres"]; then \
    COPY --from=builder /app/bin/migrate /app/bin/migrate; \
    fi

COPY . .

CMD if ["$STORAGE" = "postgres"]; then \
    /bin/sh -c "/app/bin/migrate --migrate=up && /app/bin/server"; \
    else \
    /bin/sh -c "/app/bin/server"; \
    fi

# ENTRYPOINT ["/app/bin/migrate", "--migrate=up"]
# CMD ["/app/bin/server"]