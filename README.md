# Simple Go Blog API

Built on top of Golang with the help of `go-chi`, `go-jwt`, `mysql`, `aws-sdk-go`

API specification can be found at `Insomnia.json`

Database Schema can be found at `blog.sql`

## Configure
`cp .env.example .env && vi .env`

## Development
With docker: `docker-compose up --build -d`

Without docker: `go run main.go`

Configure in `docker-compose.yml` file, default running on port 80

## Production
Build image: `sudo docker build -t go-blog .`

Run container: `sudo docker run -dit -p 80:80 --name go-blog go-blog:latest`

Kill and remove: `(sudo docker kill go-blog || true) && sudo docker rm go-blog`

Without docker: `go build && ./go-blog`