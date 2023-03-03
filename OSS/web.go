package main

import (
	"embed"
	"fileUpload/model"
	"fmt"
	"io/ioutil"
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
	r.POST("/multiUp", multiUpload)

	r.Run(":80")
}

func list(c *gin.Context) {
	objects := model.ListOb("samli007")
	c.JSON(http.StatusOK, gin.H{"objects": objects})
}

func down(c *gin.Context) {
	name := c.Param("id")
	model.DownOb("samli007", name)
	c.File("files/" + name)
}

func del(c *gin.Context) {
	if strings.Contains(c.Request.URL.Path, "code") {
		name := c.Param("id")
		err := model.DeleteOb("samli007", "code/"+name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "delete object success."})
		}
	} else {
		name := c.Param("id")
		fmt.Println(name)
		err := model.DeleteOb("samli007", name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "delete object success."})
		}
	}
}

func upload(c *gin.Context) {
	file, err := c.FormFile("f1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	fileHandle, err := file.Open() //打开上传文件
	if err != nil {
		return
	}
	defer fileHandle.Close()
	fileByte, err := ioutil.ReadAll(fileHandle) //获取上传文件字节流
	if err != nil {
		return
	}

	model.UpOb("samli007", file.Filename, fileByte)
	objects := model.ListOb("samli007")
	c.JSON(http.StatusOK, gin.H{"objects": objects})
}

func multiUpload(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["f2"]

	for _, file := range files {
		fileHandle, err := file.Open() //打开上传文件
		if err != nil {
			return
		}
		defer fileHandle.Close()

		fileByte, err := ioutil.ReadAll(fileHandle) //获取上传文件字节流
		if err != nil {
			return
		}
		model.UpOb("samli007", "code/"+file.Filename, fileByte)
	}
	objects := model.ListOb("samli007")
	c.JSON(http.StatusOK, gin.H{"objects": objects})
}
