#FROM golang:1.21-alpine
#
#WORKDIR /usr/src/subscriber
#
#COPY go.mod go.sum ./
#
#RUN go mod download && go mod verify
#
#COPY . .
#RUN mkdir -p /usr/local/bin/
#RUN go mod tidy
#RUN go build -o /usr/local/bin/subscriber ./cmd/subscriber
#
#EXPOSE 8080
#
#CMD ["/usr/local/bin/subscriber"]
