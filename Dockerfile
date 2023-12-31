FROM golang

WORKDIR /app

COPY . .

RUN go build -o server

EXPOSE 9000

CMD ["./server"]