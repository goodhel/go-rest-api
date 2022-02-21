package routes

import (
	"example/web-service-gin/pkg/album"

	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	router := gin.Default()

	router.GET("/albums", album.GetAlbums)
	router.GET("/albums/artist/:name", album.AlbumsByArtist)
	router.GET("/albums/:id", album.GetAlbumById)
	router.POST("/albums", album.PostAlbum)

	return router
}
