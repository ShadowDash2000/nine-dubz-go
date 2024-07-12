package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"nine-dubz/app/controller"
	"nine-dubz/app/model"
	"os"
)

func main() {
	appIp, ok := os.LookupEnv("APP_IP")
	if !ok {
		appIp = "localhost"
	}
	appPort, ok := os.LookupEnv("APP_PORT")
	if !ok {
		appPort = "8080"
	}
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		dbHost = "localhost"
	}
	dbLogin, ok := os.LookupEnv("DB_LOGIN")
	if !ok {
		dbLogin = "root"
	}
	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		dbPassword = ""
	}
	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		dbName = "nine-dubz"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/?charset=utf8mb4&parseTime=True&loc=Local", dbLogin, dbPassword, dbHost)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	_ = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", dbName))
	dsn = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbLogin, dbPassword, dbHost, dbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
		Logger:      logger.Default.LogMode(logger.Silent),
		PrepareStmt: true,
	})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(
		&model.Movie{},
		&model.User{},
		&model.Role{},
		&model.ApiMethod{},
		&model.Token{},
		&model.File{},
	)

	routerController := controller.NewRouterController(*db)
	router := routerController.HandleRoute()

	err = http.ListenAndServe(appIp+":"+appPort, router)
	if err != nil {
		return
	}
}
