FROM golang:1.25-alpine AS build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ratelimiter ./cmd/ratelimiter/main.go

FROM scratch
WORKDIR /app

# Copia os certificados CA do sistema Alpine da fase de build, pq a imagem scratch não tem
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/ratelimiter .

EXPOSE 8080
CMD ["./ratelimiter"]