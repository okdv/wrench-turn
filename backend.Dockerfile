FROM golang:1.20-alpine AS builder 

WORKDIR /api 

COPY go.* .
RUN go mod download 
COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./app ./main.go

FROM alpine:latest  

WORKDIR /api 
COPY --from=builder /api/app .
EXPOSE 8080 
ENTRYPOINT ["./app"]