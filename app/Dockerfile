FROM golang:1.24.2

ARG EMAIL_GITHUB
ARG NAME_GITHUB

RUN git config --global user.email ${EMAIL_GITHUB}
RUN git config --global user.name ${NAME_GITHUB}

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build .

EXPOSE 5000

CMD [ "go", "run", "." ]
