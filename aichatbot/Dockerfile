FROM golang:1.23-alpine AS build-stage
WORKDIR /aichatbot
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env .
RUN go build -o backend-api .
FROM alpine:latest
COPY --from=build-stage /aichatbot/backend-api .
COPY --from=build-stage /aichatbot/.env .
EXPOSE 8080
CMD ["./backend-api"]