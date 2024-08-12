FROM golang:alpine

RUN mkdir /publisher
WORKDIR /publisher

COPY ../../../go.mod .
COPY ../../../go.sum .
RUN go mod download

COPY . .

RUN go build -o publisher ./server/order/nats_prod/main.go

EXPOSE 8080

CMD [ "./publisher" ]