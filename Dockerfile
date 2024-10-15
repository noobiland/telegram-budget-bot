FROM golang:1.23

WORKDIR /usr/src/app

COPY go.mod go.sum ./

COPY . .

RUN go build -v -o /usr/local/bin/budget-bot ./...

CMD [ "budget-bot" ]