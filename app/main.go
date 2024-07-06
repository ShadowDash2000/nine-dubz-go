package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/controller"
	"nine-dubz/app/model"
	"os"
)

const publicDir = "public/"

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, publicDir+"dist/index.html")
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, publicDir+"admin.html")
}

func main() {
	appIp, ok := os.LookupEnv("APP_IP")
	if !ok {
		appIp = "localhost"
	}
	appPort, ok := os.LookupEnv("APP_PORT")
	if !ok {
		appPort = "25565"
	}
	dbLogin, ok := os.LookupEnv("DB_LOGIN")
	if !ok {
		dbLogin = "root"
	}
	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		dbPassword = ""
	}

	dsn := dbLogin + ":" + dbPassword + "@tcp(localhost:3306)/nine-dubz?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(
		&model.Movie{},
		&model.User{},
		&model.Role{},
		&model.ApiMethod{},
	)

	routerController := controller.NewRouterController(*db)
	router := routerController.HandleRoute()

	err = http.ListenAndServe(appIp+":"+appPort, router)
	if err != nil {
		return
	}
}
