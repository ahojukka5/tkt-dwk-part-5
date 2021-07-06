package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var main_app_port = ":" + Getenv("MAIN_APP_PORT", "3000")
var pingpong_app_host = Getenv("PINGPONG_APP_HOST", "http://pingpong-svc:8000")
var welcome_message = Getenv("MESSAGE", "Main app message not defined!")
var random_string = "undefined"

func Healthz(w http.ResponseWriter, r *http.Request) {
	// we are happy if we get _any_ response from pingpong/healthz
	_, err := http.Get(pingpong_app_host + "/healthz")
	if err != nil {
		fmt.Println("Main app health check: not ready")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		fmt.Println("Main app health check: ready")
		w.WriteHeader(http.StatusOK)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := ioutil.ReadFile("/cache/timestamp")
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(pingpong_app_host)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		msg := fmt.Sprintf(`{"message":"%s"}`, err.Error())
		fmt.Println("Cannot fetch data from ping/pong service")
		fmt.Println(msg)
		fmt.Println(err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	} else {
		fmt.Println("Ping/pong response from", pingpong_app_host)
	}
	fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
	if resp.StatusCode != 200 {
		msg := "Didn't get status code 200 from ping/pong service"
		http.Error(w, msg, http.StatusInternalServerError)
		fmt.Println(msg)
		return
	}
	fmt.Fprintln(w, welcome_message)
	fmt.Fprintf(w, "%s: %s.\n", t, random_string)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	cnt := string(body)
	fmt.Fprintf(w, "Ping Pongs: %s\n", cnt)
}

func main() {
	random_string = uuid.New().String()
	fmt.Println("Random string =", random_string)
	http.HandleFunc("/", index)
	http.HandleFunc("/healthz", Healthz)
	host := "http://localhost" + main_app_port
	println("Main app server listening in address", host)
	err := http.ListenAndServe(main_app_port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
