package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

func main() {

	var wg sync.WaitGroup
	devs := dev("dev.txt")
	lines := len(devs)
	wg.Add(lines)

	for _, line := range devs {
		go session(line, &wg)
	}
	wg.Wait()
	//time.Sleep(3 * time.Second)
}

func session(row string, wg *sync.WaitGroup) {
	if len(row) != 0 {
		strings.ReplaceAll(row, "\n", "")
		info := strings.Split(row, " ")
		addr := info[0]
		user := info[1]
		passwd := info[2]
		sshGo(user, passwd, addr, "cli.txt")
	}
	wg.Done()
}

func sshGo(user string, passwd string, addr string, txt string) {
	addr1 := fmt.Sprintf("%s%s", addr, ":22")
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", addr1, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	cli := coms(txt)

	res, err := session.CombinedOutput(cli)
	if err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}

	wf(addr+".log", string(res))
	log.Printf("log file: %s.log\n", addr)
}

func dev(path string) []string {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	readBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	res := strings.Split(string(readBytes), "\n")
	return res
}

func coms(path string) string {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	readBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	res := strings.ReplaceAll(string(readBytes), "\n", ";")
	return res
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