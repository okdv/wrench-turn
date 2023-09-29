FROM golang:1.20-alpine AS backend-builder 

RUN mkdir -p /api 
WORKDIR /api 

COPY go.* .
RUN go mod download 
COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./app ./main.go

FROM node:20-alpine AS frontend-build 

WORKDIR /app 
COPY ./frontend /app/
RUN npm ci 
RUN npm run build 

FROM alpine:latest  

WORKDIR /app 
COPY --from=backend-builder /api/app .
COPY --from=frontend-build /app/build .
EXPOSE 8080 
ENTRYPOINT ["./app"]