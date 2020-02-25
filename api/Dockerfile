# base image
FROM golang:latest

# set working directory
WORKDIR /app

# install and cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# add app
COPY . .

# build app
RUN go build ./cmd/api

# expose port and start
EXPOSE 5000
CMD ["./api"]