FROM golang:1.24

RUN apt update
RUN apt install postgresql-client nano -y

WORKDIR /app

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY "./cmd" "./cmd"
COPY "./internal" "./internal"
COPY "./migrations" "./migrations"


RUN GOOS=linux go build -o auth ./cmd/main.go
RUN chmod +x ./auth
CMD ["./auth"]