FROM golang:1.22.2

WORKDIR /app

COPY ./app/go.mod ./app/go.sum ./

RUN go mod download

COPY ./app .

RUN go build .

EXPOSE 5000

CMD [ "go", "run", "." ]
