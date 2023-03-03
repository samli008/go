package main

import (
	"embed"
	"fileUpload/model"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	//go:embed public
	staticFS embed.FS
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/public", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})
	r.GET("/public/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})

	r.GET("/list", list)
	r.GET("/down/:id", down)
	r.DELETE("/del/:id", del)
	r.DELETE("/del/code/:id", del)
	r.POST("/up", upload)

	r.Run(":80")
}

var bk = "sam"

func list(c *gin.Context) {
	objects := model.ListOb(bk)
	c.JSON(http.StatusOK, gin.H{"objects": objects})
}

func down(c *gin.Context) {
	name := c.Param("id")
	model.DownOb(bk, name)
	c.File(name)
}

func del(c *gin.Context) {
	if strings.Contains(c.Request.URL.Path, "code") {
		name := c.Param("id")
		err := model.DeleteOb(bk, "code/"+name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "delete object success."})
		}
	} else {
		name := c.Param("id")
		err := model.DeleteOb(bk, name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "delete object success."})
		}
	}
}

func upload(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["f1"]

	for _, file := range files {
		model.UpOb(bk, file.Filename, file, "")
	}
	c.Redirect(http.StatusMovedPermanently, "/public")
}
