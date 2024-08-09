package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"nine-dubz/internal/apimethod"
	"nine-dubz/internal/comment"
	"nine-dubz/internal/file"
	"nine-dubz/internal/googleoauth"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/role"
	"nine-dubz/internal/token"
	"nine-dubz/internal/user"
	"nine-dubz/internal/view"
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
		dbName = "nine_dubz"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbLogin, dbPassword, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
		Logger:         logger.Default.LogMode(logger.Silent),
		PrepareStmt:    true,
		TranslateError: true,
	})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(
		&file.File{},
		&role.Role{},
		&user.User{},
		&apimethod.ApiMethod{},
		&token.Token{},
		&googleoauth.AuthorizeState{},
		&movie.Movie{},
		&comment.Comment{},
		&view.View{},
	)

	var count int64
	db.Model(&role.Role{}).Where("code = ?", "all").Count(&count)
	if count == 0 {
		db.Create(&role.Role{Code: "all", Name: "all"})
	}

	db.Model(&role.Role{}).Where("code = ?", "admin").Count(&count)
	if count == 0 {
		db.Create(&role.Role{Code: "admin", Name: "admin"})
	}

	return db
}
