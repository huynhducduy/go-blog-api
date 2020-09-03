FROM golang:alpine AS develop

RUN apk add ca-certificates git
WORKDIR /root
COPY . .
RUN go mod download

FROM develop AS build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine AS production

#RUN apk add ca-certificates rsync openssh
WORKDIR /root
COPY --from=build /root/go-blog /root/go-blog
COPY .env .env

EXPOSE 80

ENTRYPOINT ["./go-blog"]

# sudo docker build -t go-blog .
# sudo docker run -dit -p 80:80 --name go-blog go-blog:latest
# (sudo docker kill go-blog || true) && sudo docker rm go-blog