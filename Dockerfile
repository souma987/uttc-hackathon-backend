# Build stage
FROM golang:1.25 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Final stage
FROM gcr.io/distroless/static-debian13
COPY --from=build /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
