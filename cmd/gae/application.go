package main

import (
	"github.com/wilsonth122/money-tracker-api/pkg/app"
)

func init() {
	app.Setup()
}

func main() {
	app.Start()
}
