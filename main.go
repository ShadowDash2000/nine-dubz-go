package main

import (
	"nine-dubz/app"
	gormDb "nine-dubz/db"
)

func main() {
	db := gormDb.NewGormDb()
	app := app.NewApp(*db)
	app.Start()
}
