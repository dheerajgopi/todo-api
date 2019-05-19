# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# golang 1.12 base image
FROM golang:1.12-alpine3.9

# install git
RUN apk update && apk add git

# maintainer
LABEL maintainer="Dheeraj Gopinath <dheerajgopinath@gmail.com>"

# set current working directory
WORKDIR /todo-api/src

# force go compiler to use go modules
ENV GO111MODULE=on

# copy everything from PWD to current working directory inside container
COPY . .

# download dependencies
RUN go mod tidy

# vendor
RUN go mod vendor

# install package
RUN go install

# copy config files
RUN mkdir $GOPATH/bin/config && cp config/dev.json $GOPATH/bin/config/ && cp config/prod.json $GOPATH/bin/config/

# expose port
EXPOSE 8080

# Run executable
CMD ["todo-api"]
