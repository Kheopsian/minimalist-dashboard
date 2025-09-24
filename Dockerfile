# --- STAGE 1 : BUILD ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 creates a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /dashboard-api

# --- STAGE 2 : FINAL ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates docker-cli

WORKDIR /app

COPY --from=builder /app/frontend ./frontend
COPY --from=builder /dashboard-api .

EXPOSE 8080

CMD [ "./dashboard-api" ]