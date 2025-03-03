package main

import "orbit-app/internal/app"

func main() {
	if err := app.New().Run(); err != nil {
		panic(err)
	}
}
