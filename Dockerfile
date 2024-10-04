FROM golang:1.23.1

COPY . /app
WORKDIR /app/cmd/

RUN go mod download

RUN go build -o ./server

EXPOSE 8000 

CMD ["./server"]