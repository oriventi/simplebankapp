#Build stage
FROM golang:1.20.3-alpine3.17 AS builder
WORKDIR /app

#Fetch dependencies
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

#Build go mod
COPY . .
RUN go build -o main main.go


#Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration
RUN chmod 777 /app/wait-for.sh
RUN chmod 777 /app/start.sh
RUN chmod 777 /app/main

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]