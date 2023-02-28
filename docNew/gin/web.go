package main

import (
	"embed"
	"liyang/models"
	"log"
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
	err := models.ConnectDatabase()
	checkErr(err)
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/public", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})
	r.GET("/public/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})

	linux := r.Group("/linux")
	{
		linux.GET("/docs", getDocs)
		linux.GET("/doc/:id", getDocById)
		linux.GET("/doc/name/:id", getDocById)
		linux.GET("/doc/content/:id", getDocByCon)
		linux.GET("/content/:id", getDocByName)
		linux.POST("/doc", addDoc)
		linux.PUT("/doc/:id", updateDoc)
		linux.DELETE("/:id", deleteDoc)
	}

	netapp := r.Group("/netapp")
	{
		netapp.GET("/docs", getDocs)
		netapp.GET("/doc/:id", getDocById)
		netapp.GET("/doc/name/:id", getDocById)
		netapp.GET("/doc/content/:id", getDocByCon)
		netapp.GET("/content/:id", getDocByName)
		netapp.POST("/doc", addDoc)
		netapp.PUT("/doc/:id", updateDoc)
		netapp.DELETE("/:id", deleteDoc)
	}

	dell := r.Group("/dell")
	{
		dell.GET("/docs", getDocs)
		dell.GET("/doc/:id", getDocById)
		dell.GET("/doc/name/:id", getDocById)
		dell.GET("/doc/content/:id", getDocByCon)
		dell.GET("/content/:id", getDocByName)
		dell.POST("/doc", addDoc)
		dell.PUT("/doc/:id", updateDoc)
		dell.DELETE("/:id", deleteDoc)
	}

	private := r.Group("/private")
	{
		private.GET("/docs", getDocs)
		private.GET("/doc/:id", getDocById)
		private.GET("/doc/name/:id", getDocById)
		private.GET("/doc/content/:id", getDocByCon)
		private.GET("/content/:id", getDocByName)
		private.POST("/doc", addDoc)
		private.PUT("/doc/:id", updateDoc)
		private.DELETE("/:id", deleteDoc)
	}

	vmware := r.Group("/vmware")
	{
		vmware.GET("/docs", getDocs)
		vmware.GET("/doc/:id", getDocById)
		vmware.GET("/doc/name/:id", getDocById)
		vmware.GET("/doc/content/:id", getDocByCon)
		vmware.GET("/content/:id", getDocByName)
		vmware.POST("/doc", addDoc)
		vmware.PUT("/doc/:id", updateDoc)
		vmware.DELETE("/:id", deleteDoc)
	}

	windows := r.Group("/windows")
	{
		windows.GET("/docs", getDocs)
		windows.GET("/doc/:id", getDocById)
		windows.GET("/doc/name/:id", getDocById)
		windows.GET("/doc/content/:id", getDocByCon)
		windows.GET("/content/:id", getDocByName)
		windows.POST("/doc", addDoc)
		windows.PUT("/doc/:id", updateDoc)
		windows.DELETE("/:id", deleteDoc)
	}

	r.Run(":80")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getDocs(c *gin.Context) {
	path := strings.Split(c.Request.URL.Path, "/")
	t := path[1]

	docs, err := models.GetDocs(t)
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
	path := strings.Split(c.Request.URL.Path, "/")
	t := path[1]
	id := c.Param("id")

	docs, err := models.GetDocByName(id, t)
	checkErr(err)
	// if the name is blank we can assume nothing is found
	if len(docs) == 0 {
		c.JSON(http.StatusOK, gin.H{"docs": docs})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"docs": docs})
	}
}

func getDocByCon(c *gin.Context) {
	path := strings.Split(c.Request.URL.Path, "/")
	t := path[1]
	id := c.Param("id")

	docs, err := models.GetDocByCon(id, t)
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
	path := strings.Split(c.Request.URL.Path, "/")
	t := path[1]

	id := c.Param("id")
	docs, err := models.GetDocByName(id, t)
	checkErr(err)
	// if the name is blank we can assume nothing is found
	if len(docs) == 0 {
		c.JSON(http.StatusOK, gin.H{"doc": docs})
		return
	} else {
		// html := models.Md(docs[0].Content)
		html := docs[0]
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
	path := strings.Split(c.Request.URL.Path, "/")
	t := path[1]
	name := c.Param("id")

	success, err := models.DeleteDoc(name, t)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}
