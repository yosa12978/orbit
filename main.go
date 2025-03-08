package main

import (
	"embed"
	"orbit-app/internal/app"
)

//go:embed templates/*
var templates embed.FS

//go:embed assets/*
var assets embed.FS

func main() {
	if err := app.New(templates, assets).Run(); err != nil {
		panic(err)
	}
}
