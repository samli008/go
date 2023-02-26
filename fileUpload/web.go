package main

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

	r.POST("/up", upload)
	r.POST("/multiUp", multiUpload)
	r.Run(":80")
}

func upload(c *gin.Context) {
	file, err := c.FormFile("f1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println(file.Filename)
	dst := fmt.Sprintf("./files/%s", file.Filename)
	c.SaveUploadedFile(file, dst)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
	})
}

func multiUpload(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["f2"]

	for _, file := range files {
		log.Println(file.Filename)
		dst := fmt.Sprintf("./files/%s", file.Filename)
		c.SaveUploadedFile(file, dst)
	}
	list("./files")

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d files uploaded!", len(files)),
	})
}

func list(path string) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, info := range fileInfos {
		fmt.Println(info.Name())
	}
}
