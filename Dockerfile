
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
FROM golang:1.17
RUN mkdir /go/src/pb-backend
WORKDIR /go/src/pb-backend
COPY . .

run export GO111MODULE=off
run export GOPATH=$HOME/go
run export PATH=$PATH:$GOPATH/bin
RUN go mod download
run go mod tidy

RUN CGO_ENABLED=0 go build

CMD ["./pb-backend"]
## our newly created binary executable

