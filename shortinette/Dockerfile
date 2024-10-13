FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go pkg/ ./

RUN go build .

EXPOSE 5000

CMD [ "go", "run", "." ]