FROM golang:1.20-alpine AS builder
COPY ./api /api
WORKDIR /api
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -o go-surl-api

FROM alpine AS runner
COPY --from=builder /api/go-surl-api .
CMD ["./go-surl-api"]