FROM golang

WORKDIR /app

COPY . .

EXPOSE 9000

CMD ["go", "run", "server.go"]