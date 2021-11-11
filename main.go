package main

import (
	"github.com/EvgeniyChernoskov/videoCatalog/log"
	"github.com/EvgeniyChernoskov/videoCatalog/pkg/controllers"
	"github.com/EvgeniyChernoskov/videoCatalog/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.InitLogger()
	InitConfig()
	initEnv()
	gin.SetMode(gin.ReleaseMode)
}

func main() {

	db, err := repository.ConnectDB()
	if err!=nil {
		log.Logger.Fatal(err)
	}

	repo := repository.New(db)
	ctrl := controllers.New(repo)

	go func() {
		r := gin.Default()
		r.GET("/videos", ctrl.GetVideos())
		r.GET("/videos/:id", ctrl.GetVideo())
		r.POST("/videos/", ctrl.AddVideo())
		r.PUT("/videos/", ctrl.UpdateVideo())
		r.DELETE("/videos/:id", ctrl.DeleteVideo())
		r.Run()
	}()

	log.Logger.Info("start app")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Logger.Info("shutdown app")
	db.Close()
	log.Logger.Sync()

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
