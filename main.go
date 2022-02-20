package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
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

var db *sql.DB

func main() {
	// Capture connection properties
	//cfg := mysql.Config{
	//	User: os.Getenv("DBUSER"),
	//	//Passwd: os.Getenv("DBPASS"),
	//	Net:    "tcp",
	//	Addr:   "127.0.0.1:3306",
	//	DBName: "recordings",
	//}
	// Get a database handle
	// "root:password1@tcp(127.0.0.1:3306)/test"
	var err error
	//db, err = sql.Open("mysql", cfg.FormatDSN())
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/recordings")

	if err != nil {
		log.Fatal(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/artist/:name", albumsByArtist)
	router.GET("/albums/:id", getAlbumById)
	router.POST("/albums", postAlbum)

	router.Run("localhost:8080")
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(c *gin.Context) {
	name := c.Param("name")
	// An albums slice to hold data from returned row
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
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
func getAlbums(c *gin.Context) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
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

func postAlbum(c *gin.Context) {
	var newAlbum CreateAlbum

	// Call BindJSON to bind the received JSON to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Insert data to database
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)",
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

func getAlbumById(c *gin.Context) {
	var album Album
	id := c.Param("id")

	// Loop over the list of albums,looking for
	// an album whose ID value matches the parameter
	err := db.QueryRow("SELECT * FROM album WHERE id = ?", id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No result"})
		return
	}

	c.JSON(http.StatusOK, album)
}
