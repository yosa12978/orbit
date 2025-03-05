package main

import (
	"embed"
	"orbit-app/internal/app"
)

//go:embed templates/*
var templates embed.FS

func main() {
	if err := app.New(templates).Run(); err != nil {
		panic(err)
	}
}
