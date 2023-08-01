package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Post struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Content string    `json:"contentHtml"`
}

var posts = []Post{
	{ID: "1", Title: "Blue Train", Date: time.Now(), Content: "This is ID1 Content This is ID1 Content This is ID1 Content This is ID1 Content This is "},
	{ID: "2", Title: "Jeru", Date: time.Now(), Content: "This is ID1 ContentThis is ID2 ContentThis is ID2 ContentThis is ID2 ContentThis is ID2 ContentThis is ID2 Content"},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Date: time.Now(), Content: "This is ID3 ContentThis is ID3 ContentThis is ID3 ContentThis is ID3 ContentThis is ID3 ContentThis is ID3 Content"},
}

func main() {
	err := godotenv.Load(".env") // envファイルのパスを渡す。何も渡さないと、どうディレクトリにある、.envファイルを探す
	if err != nil {
		panic("Error loading .env file")
	}
	hostDomain := os.Getenv("HOST_DOMAIN")

	router := gin.Default()
	router.GET("/api/v1/posts", getPosts)
	router.GET("/api/v1/post/:id", getPostByID)

	router.Run(hostDomain)
}

// getAlbums responds with the list of all albums as JSON.
func getPosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, posts)
}

func getPostByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range posts {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "post not found"})
}
