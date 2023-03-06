package models

import (
	"log"
	"fmt"
	"os"
	"time"
	"github.com/chatgp/gpt3"
)

func Gpt(question string) string {
	apiKey := "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"

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

