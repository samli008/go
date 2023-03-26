package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"golang.org/x/crypto/ssh"
)

func main() {
	var (
		username string
		password string
		addr     string
	)

	fmt.Print("pls input ip addr: ")
	fmt.Scan(&addr)

	fmt.Print("pls input username: ")
	fmt.Scan(&username)

	fmt.Print("pls input passwd: ")
	fmt.Scan(&password)

	addr = fmt.Sprintf("%s%s", addr, ":22")

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
	
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	cli := rf("cli.txt")
	
	res, err := session.CombinedOutput(cli)
	if err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}

	wf("ssh.log", string(res))
	fmt.Println("log output filename: ssh.log")
}

func wf(filename string, data string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0200)
	if err != nil {
		fmt.Println(err)
		return
	}
	content := fmt.Sprintf("%s%c", data, '\n')
	file.WriteString(content)
	file.Close()
}

func rf(path string) string {

	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
			panic(err)
	}

	defer file.Close()

	readBytes, err := ioutil.ReadAll(file)
	if err != nil {
			panic(err)
	}

	res := strings.ReplaceAll(string(readBytes), "\r\n",";")
	return res
}
