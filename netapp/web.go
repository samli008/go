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
	//r.Static("/web", "./public")
	r.GET("/public", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})
	r.GET("/public/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})

	r.GET("/person", getPersons)
	r.GET("/person/:id", getPersonById)
	r.POST("/person", addPerson)
	r.PUT("/person/:id", updatePerson)
	r.DELETE("/person/:id", deletePerson)
	r.Run(":80")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getPersons(c *gin.Context) {

	persons, err := models.GetPersons(1000)
	checkErr(err)

	if len(persons) == 0 {
		println(len(persons))
		c.JSON(http.StatusOK, gin.H{"person": "no"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"person": persons})
	}
}

func getPersonById(c *gin.Context) {

	id := c.Param("id")

	persons, err := models.GetPersonByName(id)
	checkErr(err)
	// if the name is blank we can assume nothing is found
	if len(persons) == 0 {
		//c.JSON(http.StatusBadRequest, gin.H{"person": persons})
		c.JSON(http.StatusOK, gin.H{"person": persons})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"person": persons})
	}
}

func addPerson(c *gin.Context) {

	var json models.Person

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddPerson(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func updatePerson(c *gin.Context) {

	var json models.Person

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	personId := c.Param("id")

	success, err := models.UpdatePerson(json, personId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deletePerson(c *gin.Context) {

	fso := c.Param("id")

	success, err := models.DeletePerson(fso)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}
