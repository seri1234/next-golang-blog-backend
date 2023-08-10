package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

type PostInput struct {
	Title   string `json:"title"`
	Content string `json:"contentHtml"`
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
	router.POST("/api/v1/post", postPost)

	router.Run(hostDomain)
}

func getPosts(c *gin.Context) {
	posts, err := getAllPosts()

	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "all posts not found"})
	}

	c.IndentedJSON(http.StatusOK, posts)
	fmt.Printf("allPosts found: %v\n", posts)
}

func getPostByID(c *gin.Context) {
	id := c.Param("id")

	post, err := postByID(id)
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "post not found"})
	}

	c.IndentedJSON(http.StatusOK, post)
	fmt.Printf("Post found: %v\n", post)
	return
}

func postPost(c *gin.Context) {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	hostDomain := os.Getenv("HOST_DOMAIN")

	var post PostInput
	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lastID, rowCnt, err := createPost(post)
	if err != nil {
		log.Fatalln(err)
		return
	}

	c.Header("location", hostDomain+"/api/v1/post/"+strconv.FormatInt(lastID, 10))
	c.JSON(http.StatusOK, gin.H{
		"status": "200",
		"data":   "success",
	})

	log.Printf("CreateComptration ID = %d, affected = %d \n", lastID, rowCnt)
	if err != nil {
		log.Fatalln(err)
	}

}

func postByID(id string) ([]Post, error) {
	var post []Post

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
		post = append(post, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postByID %q: %v", id, err)
	}

	return post, nil
}

func getAllPosts() ([]Post, error) {
	var posts []Post

	rows, err := db.Query("SELECT id,title,created_at FROM posts")
	if err != nil {
		return nil, fmt.Errorf("getAllPosts : %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("getAllPosts : %v", err)
		}
		posts = append(posts, p)

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("getAllPosts : %v", err)
		}
	}
	return posts, nil
}

func createPost(post PostInput) (int64, int64, error) {
	stmt, err := db.Prepare("INSERT INTO posts(title, content) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.Title, post.Content)
	if err != nil {
		log.Fatal(err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return lastID, rowCnt, err
}
