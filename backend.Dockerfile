FROM golang:1.22.2-alpine3.19 AS build 

ARG GO_ENV=production
ENV GO_ENV=production

RUN apk update && apk add --no-cache gcc musl-dev

WORKDIR /app 

COPY go.* ./
RUN go mod download 

COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o wrench-turn ./

FROM alpine:latest  

ARG GO_ENV=production
ENV GO_ENV=production

WORKDIR /app 
COPY --from=build /app/wrench-turn ./
# create data dir
RUN mkdir /app/data

EXPOSE 8080 
ENTRYPOINT ["./wrench-turn"]