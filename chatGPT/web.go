
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
	r.Run(":8181")
}

func gpt(c *gin.Context){
	question := c.Param("id")
	apiKey:= models.Rf("api.key")
	answer := models.Gpt(question,apiKey)
	models.Wf("gpt.log",question)
	c.JSON(http.StatusOK, gin.H{"chat": answer})
}