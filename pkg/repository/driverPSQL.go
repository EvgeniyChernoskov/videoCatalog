package repository

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"os"
)
type Config struct {
	DBName   string
	Host     string
	Port     string
	User     string
	SSLMode  string
	Password string
}

func ConnectDB() (*sql.DB,error) {

	config := Config{
		DBName:   viper.GetString("db.dbname"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.username"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD")}

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		config.Host, config.Port, config.User, config.DBName, config.SSLMode, config.Password)

	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil,err
	}

	if err := db.Ping(); err != nil {
		return nil,err
	}
	return db,nil
}



