FROM golang:1.23

RUN apt-get update && apt-get install -y --no-install-recommends make && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "test", "-bench=.", "-benchmem"]
