package models

import (
	"log"
	"time"
	"github.com/chatgp/gpt3"
)

func Gpt(question string) string {
	apiKey := "sk-xxx"

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

