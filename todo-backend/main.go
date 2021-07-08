package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ItemWithoutID struct {
	Task string `json:"task"`
	Done bool   `json:"done"`
}

type Item struct {
	ID   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Task string             `json:"task"`
	Done bool               `json:"done"`
}

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var port = ":" + Getenv("TODO_BACKEND_PORT", "8000")
var username = Getenv("MONGO_USERNAME", "root")
var password = Getenv("MONGO_PASSWORD", "")
var host = Getenv("MONGO_HOST", "localhost")
var mongo_uri = "mongodb://" + username + ":" + password + "@" + host
var nats_url = Getenv("NATS_URL", "nats://localhost:4222")

func GetClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(mongo_uri)
	client, err := mongo.Connect(ctx, clientOptions)
	return client, err
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	client, _ := GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if client.Ping(ctx, readpref.Primary()) == nil {
		log.Println("Todo backend app health check: ready")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Todo backend app health check: not ready")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	log.Println("getTodos")
	ctx := context.TODO()
	client, err := GetClient()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		msg := fmt.Sprintf(`{"message":"%s"}`, err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println("getTodos failed:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var items []Item
	log.Println("getTodos: fetching collection")
	collection := client.Database("todo").Collection("items")
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		msg := fmt.Sprintf(`{"message":"%s"}`, err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println("getTodos failed:", err)
		return
	}
	defer cur.Close(ctx)
	log.Println("getTodos: looping collection")
	for cur.Next(ctx) {
		var item Item
		err := cur.Decode(&item)
		if err != nil {
			fmt.Println(err)
			return
		}
		items = append(items, item)
	}
	log.Println("getTodos: encoding to json")
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func postTodo(w http.ResponseWriter, r *http.Request) {
	log.Println("postTodo")
	client, err := GetClient()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		msg := fmt.Sprintf(`{"message":"%s"}`, err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println("getTodos failed:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var item Item
	var itemWithoutID ItemWithoutID
	err = json.NewDecoder(r.Body).Decode(&itemWithoutID)
	if err != nil {
		log.Println("Failed to parse input:", err)
		return
	}
	item.ID = primitive.NewObjectID()
	item.Task = itemWithoutID.Task
	item.Done = false
	log.Println("postTodo: new todo item: " + item.Task)
	if len(item.Task) > 140 {
		log.Println("postTodo: message is too long, over 140 characters!")
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	collection := client.Database("todo").Collection("items")
	result, err := collection.InsertOne(context.Background(), item)
	if err != nil {
		log.Println("Failed to insert to database:", err)
		return
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Failed to encode to json:", err)
		return
	}
	nc, _ := nats.Connect(nats_url)
	nc.Publish("todo", []byte(`{"user": "todo-bot", "message": "New Todo item created!"}`))
	nc.Close()
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("postTodo with id %s", vars["id"])
	client, err := GetClient()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		msg := fmt.Sprintf(`{"message":"%s"}`, err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println("getTodos failed:", err)
		return
	}
	collection := client.Database("todo").Collection("items")
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"done": true}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		msg := fmt.Sprintf(`{"message":"%s"}`, err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println("getTodos failed:", err)
	}
	nc, _ := nats.Connect(nats_url)
	nc.Publish("todo", []byte(`{"user": "todo-bot", "message": "Todo item updated!"}`))
	nc.Close()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todos", postTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	println("Server listening in address http://localhost" + port)
	nc, _ := nats.Connect(nats_url)
	nc.Publish("todo", []byte(`{"user": "todo-bot", "message": "Starting backend app!"}`))
	nc.Close()
	http.ListenAndServe(port, router)
}
