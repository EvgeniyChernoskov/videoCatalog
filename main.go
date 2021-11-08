package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Video struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

var videos []Video

func main() {

	videos = append(videos,
		Video{Id: 1, Title: "video1", Description: "desription1", Url: "htttp//www.video1.com"},
		Video{Id: 2, Title: "video2", Description: "desription2", Url: "htttp//www.video2.com"},
		Video{Id: 3, Title: "video1", Description: "desription3", Url: "htttp//www.video3.com"},
		Video{Id: 4, Title: "video1", Description: "desription4", Url: "htttp//www.video4.com"})

	r := gin.Default()

	r.GET("/videos", getVideos)
	r.GET("/videos/:id", getVideo)
	r.POST("/videos/", addVideo)
	r.PUT("/videos/", updateVideo)
	r.DELETE("/videos/:id", deleteVideo)
}

func getVideos(c *gin.Context) {
	c.JSON(http.StatusOK, videos)
}

func getVideo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	c.JSON(http.StatusOK, videos[id-1])
}

func addVideo(c *gin.Context) {
	video := Video{}
	c.BindJSON(&video)
	videos = append(videos, video)
}

func updateVideo(c *gin.Context) {
	video := Video{}
	c.BindJSON(&video)
	for i, vid := range videos {
		if video.Id == vid.Id {
			videos[i] = video
		}
	}
}

func deleteVideo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, vid := range videos {
		if id == vid.Id {
			videos = append(videos[:i], videos[i+1:]...)
		}
	}
}
