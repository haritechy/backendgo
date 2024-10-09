FROM golang:1.23-alpine
WORKDIR /aichatbot
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o backend-api .
EXPOSE 8080
CMD ["./backend-api"]