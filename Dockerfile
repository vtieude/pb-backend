
# # syntax=docker/dockerfile:1

# ##
# ## Build
# ##
# FROM golang:1.16-alpine

# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

# COPY . .

# RUN go build -o /docker-gs-ping

# EXPOSE 3000

# CMD [ "/docker-gs-ping" ]

# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17.5-alpine3.15 AS builder
RUN mkdir /go/src/pb-backend
WORKDIR /go/src/pb-backend

COPY go.mod go.sum ./
# run export GO111MODULE=off
# run export GOPATH=$HOME/go
# run export PATH=$PATH:$GOPATH/bin
RUN go mod download
# run go mod tidy
# run go mod vendor
COPY . .

RUN CGO_ENABLED=0 go build

FROM alpine
RUN adduser -S -D -H -h /go/src/app appuser
USER appuser

COPY --from=builder /go/src/pb-backend/pb-backend /app/
COPY --from=builder /go/src/pb-backend/config.yml /app/
WORKDIR /app
CMD ["./pb-backend"]



