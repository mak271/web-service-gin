package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1111"
	dbname   = "recordings"
)

var db *sql.DB

// Album represents data about a record Album.
type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// defAlbums slice to seed record album data.
var defAlbums = []Album{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrain", Price: 56.99},
	{ID: 2, Title: "Jeru", Artist: "Gerry Milligan", Price: 17.99},
	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlconn)
	checkError(err)

	defer db.Close()

	//rows, err := db.Query(`SELECT * FROM "album"`)

	// defer rows.Close()

	// for rows.Next() {
	// 	var name string
	// 	var roll string

	// 	err = rows.Scan(&name, &roll)
	// 	checkError(err)

	// 	fmt.Println(name, roll)
	// }

	checkError(err)

	// deleteStmt := `delete from "Students" where id = $1`
	// _, e := db.Exec(deleteStmt, 1)
	// checkError(e)

	// updateStmt := `update "Students" set "Name" = $1, "Roll" = $2 where "id" = $3`
	// _, e := db.Exec(updateStmt, "Mary", "3", 2)
	// checkError(e)

	// insertStmt := `insert into "Students"("Name", "Roll") values ('John', '1')`
	// _, e := db.Exec(insertStmt)
	// checkError(e)

	// insertDynStmt := `insert into "Students"("Name", "Roll") values ($1, $2)`
	// _, e := db.Exec(insertDynStmt, "Jane", "5")
	// checkError(e)

	err = db.Ping()
	checkError(err)

	fmt.Println("Connected")

	router := gin.Default()
	router.GET("/albums/:artist", getAlbumsByArtist)
	router.GET("/albums", getAlbums)
	// router.GET("/albums/:id", getAlbunByID)
	router.POST("/albums", postAlbums)
	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album")
	checkError(err)

	defer rows.Close()

	for rows.Next() {
		var alb Album
		err = rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
		checkError(err)
		albums = append(albums, alb)
	}

	if len(albums) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "albums not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum Album

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	_, err := db.Exec("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	checkError(err)

	c.IndentedJSON(http.StatusOK, gin.H{"status": true})
}

func getAlbumsByArtist(c *gin.Context) {
	artist := c.Param("artist")
	artist = addSpace(artist)

	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", artist)

	checkError(err)

	defer rows.Close()

	for rows.Next() {
		var alb Album
		err = rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
		checkError(err)
		albums = append(albums, alb)
	}

	if len(albums) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, albums)

}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
// func getAlbunByID(c *gin.Context) {
// 	id := c.Param("id")

// 	// Loop over the list of albums, looking for
// 	// an album whose ID value matches the parameter.
// 	for _, a := range albums {
// 		if a.ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
// }

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func addSpace(s string) string {
	buf := &bytes.Buffer{}
	for i, rune := range s {
		if unicode.IsUpper(rune) && i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteRune(rune)
	}
	return buf.String()
}
