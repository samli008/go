package main

import (
	"embed"
	"liyang/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed public
	staticFS embed.FS
)

func main() {
	err := models.ConnectDatabase()
	checkErr(err)

	r := gin.Default()
	r.GET("/public", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})
	r.GET("/public/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})

	r.GET("/doc", getDocs)
	r.GET("/doc/:id", getDocById)
	r.POST("/doc", addDoc)
	r.PUT("/doc/:id", updateDoc)
	r.DELETE("/doc/:id", deleteDoc)
	r.GET("/content/:id", getDocByName)
	r.Run(":80")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getDocs(c *gin.Context) {

	docs, err := models.GetDocs()
	checkErr(err)

	if len(docs) == 0 {
		println(len(docs))
		c.JSON(http.StatusOK, gin.H{"docs": "no"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"docs": docs})
	}
}

func getDocById(c *gin.Context) {

	id := c.Param("id")

	docs, err := models.GetDocByName(id)
	checkErr(err)
	// if the name is blank we can assume nothing is found
	if len(docs) == 0 {
		c.JSON(http.StatusOK, gin.H{"docs": docs})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"docs": docs})
	}
}

func getDocByName(c *gin.Context) {

	id := c.Param("id")
	docs, err := models.GetDocByName(id)
	checkErr(err)
	// if the name is blank we can assume nothing is found
	if len(docs) == 0 {
		c.JSON(http.StatusOK, gin.H{"doc": docs})
		return
	} else {
		html := models.Md(docs[0].Content)
		c.JSON(http.StatusOK, gin.H{"doc": html})
	}
}

func addDoc(c *gin.Context) {

	var json models.Doc

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddDoc(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func updateDoc(c *gin.Context) {

	var json models.Doc

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docId := c.Param("id")

	success, err := models.UpdateDoc(json, docId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deleteDoc(c *gin.Context) {

	name := c.Param("id")

	success, err := models.DeleteDoc(name)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}
