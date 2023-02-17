package main

import (
	"embed"
	"liyang/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type docker1 struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Network string `json:"network"`
	Ip      string `json:"ip"`
}

type docker2 struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Network string `json:"network"`
	Ip      string `json:"ip"`
	Count   string `json:"count"`
}

var (
	//go:embed public
	staticFS embed.FS
)

func main() {
	r := gin.Default()
	r.GET("/public", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})
	r.GET("/public/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticFS))
	})

	r.GET("/list", list)
	r.DELETE("/del/:id", del)
	r.POST("/create", create)
	r.POST("/createAuto", createAuto)
	r.GET("/images", images)
	r.GET("/listNet", listNet)
	r.Run(":80")

}

func list(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	lists := model.ContainerList(cli)

	if len(lists) == 0 {
		println(len(lists))
		c.JSON(http.StatusOK, gin.H{"lists": "no"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"lists": lists})
	}
}

func images(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	lists := model.ListImage(cli)

	if len(lists) == 0 {
		println(len(lists))
		c.JSON(http.StatusOK, gin.H{"lists": "no"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"lists": lists})
	}
}

func listNet(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	lists := model.NetList(cli)

	if len(lists) == 0 {
		println(len(lists))
		c.JSON(http.StatusOK, gin.H{"lists": "no"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"lists": lists})
	}
}

func create(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	var json docker1

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := model.CreateContainer(cli, json.Name, json.Network, json.Ip, json.Image, json.Name)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": result})
	}
}

func createAuto(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	var json docker2

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	count, _ := strconv.Atoi(json.Count)
	for i := 0; i < count; i++ {
		name := json.Name + strconv.Itoa(i)
		str1 := strings.Split(json.Ip, ".")
		num, _ := strconv.Atoi(str1[len(str1)-1])
		num += i
		str2 := str1[:len(str1)-1]
		str2 = append(str2, strconv.Itoa(num))
		ip := str2[0] + "." + str2[1] + "." + str2[2] + "." + str2[3]

		_, err := model.CreateContainer(cli, name, json.Network, ip, json.Image, name)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": "created fail"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "created success"})
}

func del(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	id1 := c.Param("id")
	_, err = model.RemoveContainer(id1, cli)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": " delete failed."})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": " delete success."})
	}
}
