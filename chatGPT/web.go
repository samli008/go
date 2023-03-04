
package main

import (
	"embed"
	"liyang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed public
	staticFS embed.FS
)

func main() {
	r := gin.Default()
	//r.Static("/web", "./public")
	r.GET("/public", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})
	r.GET("/public/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})

	r.POST("/gpt/:id", gpt)
	r.Run(":8080")
}

func gpt(c *gin.Context){
	question := c.Param("id")
	answer := models.Gpt(question)
	c.JSON(http.StatusOK, gin.H{"chat": answer})
}