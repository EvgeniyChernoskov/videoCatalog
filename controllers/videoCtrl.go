package controllers

import (
	"errors"
	"github.com/EvgeniyChernoskov/videoCatalog/log"
	"github.com/EvgeniyChernoskov/videoCatalog/models"
	"github.com/EvgeniyChernoskov/videoCatalog/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
	r repository.VideoRepository
}

func New(rep repository.VideoRepository) Controller {
	return Controller{r: rep}
}

func (ctrl Controller) GetVideos() func(c *gin.Context) {
	return func(c *gin.Context) {
		videos, err := ctrl.r.GetVideos()
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, videos)
	}
}

func (ctrl Controller) GetVideo() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			errorResponse(c, http.StatusBadRequest, errors.New("wrong id"))
			return
		}

		video, err := ctrl.r.GetVideo(id)
		if err != nil {
			errorResponse(c, http.StatusOK, errors.New("not found"))
			return
		}

		c.JSON(http.StatusOK, video)
	}
}

func (ctrl Controller) AddVideo() func(c *gin.Context) {
	return func(c *gin.Context) {
		var video models.Video

		err := c.BindJSON(&video)
		if err != nil {
			errorResponse(c, http.StatusBadRequest, errors.New("error JSON body"))
			return
		}

		videoId, err := ctrl.r.AddVideo(video)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, map[string]int{"id": videoId})
	}
}

func (ctrl Controller) UpdateVideo() func(c *gin.Context) {
	return func(c *gin.Context) {
		var video models.Video

		err := c.BindJSON(&video)
		if err != nil {
			errorResponse(c, http.StatusBadRequest, err)
			return
		}

		rowsUpdated, err := ctrl.r.UpdateVideo(video)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, map[string]int64{"rows updated": rowsUpdated})
	}
}

func (ctrl Controller) DeleteVideo() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			errorResponse(c, http.StatusBadRequest, errors.New("wrong id"))
			return
		}

		rowsDeleted, err := ctrl.r.RemoveVideo(id)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, map[string]int64{"deleted rows:": rowsDeleted})
	}
}



func errorResponse(c *gin.Context, code int, err error) {
	log.Logger.Error(err.Error())
	c.AbortWithStatusJSON(code, err.Error())
}
