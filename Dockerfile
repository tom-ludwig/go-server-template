FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

FROM gcr.io/distroless/static-debian13:nonroot 
#:debug

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
