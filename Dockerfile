FROM golang:1.24-alpine
WORKDIR /app
COPY . .
WORKDIR /app/
CMD ["go", "run", "cmd/main.go"]
ARG APP_PORT=8080
ENV APP_PORT=${APP_PORT}
EXPOSE ${APP_PORT}