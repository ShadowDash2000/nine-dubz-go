package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"nine-dubz/internal/apimethod"
	"nine-dubz/internal/file"
	"nine-dubz/internal/googleoauth"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/role"
	"nine-dubz/internal/token"
	"nine-dubz/internal/user"
	"os"
)

func NewGormDb() *gorm.DB {
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

	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", dbName))
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
		&movie.Movie{},
		&user.User{},
		&role.Role{},
		&apimethod.ApiMethod{},
		&token.Token{},
		&file.File{},
		&googleoauth.AuthorizeState{},
	)

	return db
}
