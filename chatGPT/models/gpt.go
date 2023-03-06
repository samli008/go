package models

import (
	"log"
	"fmt"
	"os"
	"time"
	"strings"
	"io/ioutil"
	"github.com/chatgp/gpt3"
)

func Gpt(question string, key string) string {
	apiKey := key

	// new gpt-3 client
	cli, _ := gpt3.NewClient(&gpt3.Options{
		ApiKey:  apiKey,
		Timeout: 30 * time.Second,
		Debug:   true,
	})

	// request api
	uri := "/v1/chat/completions"
	params := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]interface{}{
			{"role": "user", "content": question},
		},
	}

	res, err := cli.Post(uri, params)
	if err != nil {
		log.Fatalf("request api failed: %v", err)
	}

	// fmt.Printf("chatGPT answer your question: %s\n", res.Get("choices.0.message.content").String())
	answer := res.Get("choices.0.message.content").String()
	return answer
}

func Wf(filename string, data string) {
	file, err := os.OpenFile(filename, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0200)
	if err != nil {
		fmt.Println(err)
		return
	}
	content:=fmt.Sprintf("%s%c",data,'\n')
	file.WriteString(content)
	file.Close()
}

func Rf(path string) string {
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
	key := results[0]
	return key
}

