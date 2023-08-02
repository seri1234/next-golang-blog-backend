package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Content   string    `json:"contentHtml"`
}

var posts = []Post{
	{ID: "1", Title: "Blue Train", CreatedAt: time.Now(), Content: "This is ID1 Content This is ID1 Content This is ID1 Content This is ID1 Content This is "},
	{ID: "2", Title: "Jeru", CreatedAt: time.Now().AddDate(0, 0, 1), Content: "This is ID1 ContentThis is ID2 ContentThis is ID2 ContentThis is ID2 ContentThis is ID2 ContentThis is ID2 Content"},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", CreatedAt: time.Now().AddDate(0, 0, 2), Content: "This is ID3 ContentThis is ID3 ContentThis is ID3 ContentThis is ID3 ContentThis is ID3 ContentThis is ID3 Content"},
}

var db *sql.DB

func main() {
	var err error

	err = godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	hostDomain := os.Getenv("HOST_DOMAIN")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbAddr := os.Getenv("DBADDR")

	locale, _ := time.LoadLocation("Asia/Tokyo")

	cfg := mysql.Config{
		User:      dbUser,
		Passwd:    dbPass,
		Net:       "tcp",
		Addr:      dbAddr,
		DBName:    "golang_next_blogs",
		ParseTime: true,
		Loc:       locale,
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.GET("/api/v1/posts", getPosts)
	router.GET("/api/v1/post/:id", getPostByID)

	router.Run(hostDomain)
}

func getPosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, posts)
}

func getPostByID(c *gin.Context) {
	id := c.Param("id")

	post, err := postByID(id)
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "post not found"})
	}

	c.IndentedJSON(http.StatusOK, post)
	fmt.Printf("Posts found: %v\n", post)
	return
}

func postByID(id string) ([]Post, error) {
	var posts []Post

	rows, err := db.Query("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("postByID %q: %v", id, err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("postByID %q: %v", id, err)
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postByID %q: %v", id, err)
	}
	return posts, nil
}
