FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod edit -go=1.23

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata ffmpeg

ENV TZ=Asia/Jakarta

WORKDIR /root/

COPY --from=builder /app/server .

RUN mkdir -p uploads

EXPOSE 8082

CMD ["./server"]