package repository

import (
	"database/sql"
	"github.com/EvgeniyChernoskov/videoCatalog/pkg/models"
)


type VideoRepository struct {
	db *sql.DB
}
func New(db *sql.DB) VideoRepository {
	return VideoRepository{db: db}
}

func (r VideoRepository) GetVideos() ([]models.Video, error) {
	var videos []models.Video
	var video models.Video

	rows, err := r.db.Query("SELECT * FROM videos;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&video.Id, &video.Title, &video.Description, &video.Url)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func (r VideoRepository) GetVideo(id int) (models.Video, error) {
	row := r.db.QueryRow("SELECT * FROM videos WHERE id=$1;", id)
	var video models.Video
	err := row.Scan(&video.Id, &video.Title, &video.Description, &video.Url)
	if err != nil {
		return models.Video{}, err
	}
	return video, nil
}

func (r VideoRepository) AddVideo(video models.Video) (int,error) {
	var videoId int
	err := r.db.QueryRow("INSERT INTO videos (title,description,url) VALUES ($1,$2,$3) RETURNING id;",
		video.Title, video.Description, video.Url).Scan(&videoId)
	if err != nil {
		return 0, err
	}
	return videoId, nil
}

func (r VideoRepository) UpdateVideo( video models.Video) (int64, error) {
	result, err := r.db.Exec("UPDATE videos SET title=$1 , description=$2, url=$3 WHERE id=$4 RETURNING id;",
		video.Title, video.Description, video.Url, video.Id)
	if err != nil {
		return 0, err
	}
	rowsUpdated, err := result.RowsAffected()
	return rowsUpdated, nil
}

func (r VideoRepository) RemoveVideo(id int) (int64, error) {
	result, err := r.db.Exec("DELETE FROM videos WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	rowsDeleted, err := result.RowsAffected()
	return rowsDeleted, nil
}
