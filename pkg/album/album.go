package album

import (
	"example/web-service-gin/pkg/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

type CreateAlbum struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

// albumsByArtist queries for albums that have the specified artist name.
func AlbumsByArtist(c *gin.Context) {
	name := c.Param("name")
	// An albums slice to hold data from returned row
	var albums []Album

	rows, err := database.DB.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "albumsByArtist not found"})
		return
	}

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "albumsByArtist not found"})
			return
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "albumsByArtist not found"})
		return
	}

	c.JSON(http.StatusOK, albums)
}

// getAlbums responds with the list of all albums as JSON.
func GetAlbums(c *gin.Context) {
	var albums []Album

	rows, err := database.DB.Query("SELECT * FROM album")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not Found"})
			return
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not Found"})
		return
	}

	c.JSON(http.StatusOK, albums)
}

func PostAlbum(c *gin.Context) {
	var newAlbum CreateAlbum

	// Call BindJSON to bind the received JSON to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Insert data to database
	result, err := database.DB.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)",
		newAlbum.Title, newAlbum.Artist, newAlbum.Price)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error insert"})
		return
	}

	// Get the new album's generated ID for the client.
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error generated last id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": id})
}

func GetAlbumById(c *gin.Context) {
	var album Album
	id := c.Param("id")

	// Loop over the list of albums,looking for
	// an album whose ID value matches the parameter
	err := database.DB.QueryRow("SELECT * FROM album WHERE id = ?", id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No result"})
		return
	}

	c.JSON(http.StatusOK, album)
}
