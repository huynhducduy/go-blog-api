package main

import (
	log "github.com/sirupsen/logrus"
	"go-blog/cmd/app"
)

func main() {
	log.Fatal(app.Run())
}