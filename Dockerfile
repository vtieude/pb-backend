
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
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

 COPY . .

RUN go build -o /docker-base-project


##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /docker-base-project /docker-base-project

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/docker-base-project"]

