package main

import (
	"github.com/EvgeniyChernoskov/videoCatalog/controllers"
	"github.com/EvgeniyChernoskov/videoCatalog/driver"
	"github.com/EvgeniyChernoskov/videoCatalog/log"
	"github.com/EvgeniyChernoskov/videoCatalog/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func init() {
	log.InitLogger()
	InitConfig()
	initEnv()
}

func main() {
	defer log.Logger.Sync()

	db := driver.ConnectDB()
	defer db.Close()

	c :=controllers.New(repository.New(db))

	r := gin.Default()
	r.GET("/videos", c.GetVideos())
	r.GET("/videos/:id", c.GetVideo())
	r.POST("/videos/", c.AddVideo())
	r.PUT("/videos/", c.UpdateVideo())
	r.DELETE("/videos/:id", c.DeleteVideo())
	r.Run()

}

func InitConfig() {
	viper.SetConfigFile("configs/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Logger.Fatalf("error init configs: %s", err.Error())
	}
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Logger.Fatalf("Error loading .env file:%s", err.Error())
	}
}
