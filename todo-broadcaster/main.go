package main

// to test app
// 1. start todo-json-echo
// 2. ngrok http 8000
// 3. change external url to todo-broadcaster-deployment.yaml

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/nats-io/nats.go"
)

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var NATS_URL = Getenv("NATS_URL", "nats://localhost:4222")

// where to send a payload
var EXTERNAL_URL = Getenv("EXTERNAL_URL", "http://localhost:8123")

func handleQueue(msg *nats.Msg) {
	log.Println("Message from NATS: " + string(msg.Data))
	log.Println("Sending message to external service: " + EXTERNAL_URL)
	req, err := http.NewRequest("POST", EXTERNAL_URL, bytes.NewBuffer(msg.Data))
	if err != nil {
		log.Println("Failed to create HTTP POST request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send HTTP POST request: ", err)
	}
	resp.Body.Close()
	log.Println("Response Status from external service:", resp.Status)
}

func main() {
	nc, _ := nats.Connect(NATS_URL)
	nc.QueueSubscribe("todo", "queue", handleQueue)
	log.Println("Broadcasting messages from NATS to " + EXTERNAL_URL)
	runtime.Goexit()
}
