FROM golang:1.20-alpine AS builder 

ARG GO_ENV=production
ENV GO_ENV=production

WORKDIR /app 

COPY go.* .
RUN go mod download 

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o wrench-turn .

FROM alpine:latest  

ARG GO_ENV=production
ENV GO_ENV=production

WORKDIR /api 
COPY --from=builder /app .
EXPOSE 8080 
ENTRYPOINT ["./wrench-turn"]