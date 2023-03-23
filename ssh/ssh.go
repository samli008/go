package main

import (
	//"bytes"
	"fmt"
	"log"
	"io/ioutil"
	"os"
	"strings"
  
  "golang.org/x/crypto/ssh"
)

func main() {
	var (
		username string
		password string
		addr string
	)
	
	fmt.Print("pls input username: ")
	fmt.Scan(&username)

	fmt.Print("pls input passwd: ")
	fmt.Scan(&password)

	fmt.Print("pls input ip addr: ")
	fmt.Scan(&addr)

	addr =fmt.Sprintf("%s%s",addr,":22")

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	// 开启一个session，用于执行一个命令
	

	// 执行命令，并将执行的结果写到 b 中
	// var b bytes.Buffer
	// session.Stdout = &b
  
	cli := rf("cli")
	fmt.Println(cli)
	var cli2 string
	for _,v := range cli {
		cli2 = cli2 + v
	}
	fmt.Println(cli2)
	// if err := session.Run(cli); err != nil {
	// 	log.Fatal("Failed to run: " + err.Error())
	// }
	// fmt.Println(b.String())

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()
	
	res, err := session.CombinedOutput(cli2)
	if err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}

	wf("u.log",string(res))
}

func wf(filename string, data string) {
	file, err := os.OpenFile(filename, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0200)
	if err != nil {
		fmt.Println(err)
		return
	}
	content:=fmt.Sprintf("%s%c",data,'\n')
	file.WriteString(content)
	file.Close()
}

func rf(path string) []string {
	fileHanle, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer fileHanle.Close()

	readBytes, err := ioutil.ReadAll(fileHanle)
	if err != nil {
		panic(err)
	}

	results := strings.Split(string(readBytes), "\n")
	//res := string(readBytes)
	return results
}