package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"strconv"
)

type Video struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

var db *sql.DB
var logger *zap.SugaredLogger

func main() {

	InitLogger()
	defer logger.Sync()

	if err := initConfig(); err != nil {
		logger.Fatalf("error init configs: %s", err.Error())
	}
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file:%s", err.Error())
	}

	type Config struct {
		DBName   string
		Host     string
		Port     string
		User     string
		SSLMode  string
		Password string
	}
	config := Config{
		DBName:   viper.GetString("db.dbname"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.username"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD")}

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		config.Host, config.Port, config.User, config.DBName, config.SSLMode, config.Password)
	db, err = sql.Open("postgres", psqlConn)

	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal(err)
	}

	r := gin.Default()
	r.GET("/videos", getVideos)
	r.GET("/videos/:id", getVideo)
	r.POST("/videos/", addVideo)
	r.PUT("/videos/", updateVideo)
	r.DELETE("/videos/:id", deleteVideo)

	r.Run()

}

func getVideos(c *gin.Context) {
	var video Video
	var videos []Video

	rows, err := db.Query("SELECT * FROM videos;")
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&video.Id, &video.Title, &video.Description, &video.Url)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err)
			return
		}
		videos = append(videos, video)
	}
	c.JSON(http.StatusOK, videos)
}

func getVideo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, errors.New("wrong id"))
		return
	}

	row := db.QueryRow("SELECT * FROM videos WHERE id=$1;", id)

	var video Video
	err = row.Scan(&video.Id, &video.Title, &video.Description, &video.Url)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, video)
}

func addVideo(c *gin.Context) {
	var video Video
	var videoId int

	err := c.BindJSON(&video)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, errors.New("error JSON body"))
		return
	}

	err = db.QueryRow("INSERT INTO videos (title,description,url) VALUES ($1,$2,$3) RETURNING id;",
		video.Title, video.Description, video.Url).Scan(&videoId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, map[string]int{"id": videoId})
}

func updateVideo(c *gin.Context) {
	var video Video
	err := c.BindJSON(&video)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	result, err := db.Exec("UPDATE videos SET title=$1 , description=$2, url=$3 WHERE id=$4 RETURNING id;",
		video.Title, video.Description, video.Url, video.Id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, map[string]int64{"rows updated": rowsUpdated})
}

func deleteVideo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, errors.New("wrong id"))
		return
	}

	result, err := db.Exec("DELETE FROM videos WHERE id=$1", id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]int64{"deleted rows:": rowsDeleted})
}

func initConfig() error {
	viper.SetConfigFile("configs/config.yml")
	return viper.ReadInConfig()
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger = zap.New(core, zap.AddCaller()).Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
	//return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./video.log",
		MaxSize:    5,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
	//file, _ := os.Create("./test.log")
	//return zapcore.AddSync(file)
}

func errorResponse(c *gin.Context, code int, err error) {
	logger.Error(err.Error())
	c.AbortWithStatusJSON(code, err.Error())
}
