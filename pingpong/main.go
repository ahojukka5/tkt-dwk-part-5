package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var port = ":" + Getenv("PINGPONG_APP_PORT", "8000")
var pg_host = Getenv("POSTGRES_HOST", "pingpong-database-svc")
var pg_port = Getenv("POSTGRES_PORT", "5432")
var pg_user = Getenv("POSTGRES_USER", "postgres")
var pg_password = Getenv("POSTGRES_PASSWORD", "")
var pg_dbname = Getenv("POSTGRES_DATABASE", "postgres")
var db_initialized = false

func OpenConnection() (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		pg_host, pg_port, pg_user, pg_password, pg_dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB() error {
	if db_initialized {
		return nil
	}
	db, err := OpenConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS PingPongCounter (cnt int);")
	if err != nil {
		return err
	}
	var numberOfRows int
	err = db.QueryRow("SELECT COUNT(*) FROM PingPongCounter;").Scan(&numberOfRows)
	if err != nil {
		return err
	}
	log.Println("Number of rows in database =", numberOfRows)
	if numberOfRows == 0 {
		_, err = db.Exec("INSERT INTO PingPongCounter VALUES (0)")
		if err != nil {
			return err
		}
	}
	db_initialized = true
	return nil
}

func IncreasePingPongCounter() (int, error) {
	db, err := OpenConnection()
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var cnt int
	err = db.QueryRow("SELECT cnt FROM PingPongCounter LIMIT 1;").Scan(&cnt)
	if err != nil {
		return 0, err
	}
	cnt += 1
	_, err = db.Exec("UPDATE PingPongCounter SET cnt = $1;", cnt)
	if err != nil {
		return 0, err
	}
	log.Println("Updated PingPongCounter to", cnt)
	return cnt, nil
}

func RegisterPing(w http.ResponseWriter, r *http.Request) {
	cnt, err := IncreasePingPongCounter()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		msg := fmt.Sprintf(`{"message":"%s"}`, err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println("Cannot ping:", err)
	} else {
		log.Println("PingPong called, cnt =", cnt)
		fmt.Fprint(w, cnt)
	}
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	err := InitDB()
	if err != nil {
		log.Println("Unable to initialize database: ", err)
	}
	if db_initialized {
		log.Println("Ping/pong app health check: ready")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Ping/pong app health check: not ready")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	err := InitDB()
	if err != nil {
		log.Println("Unable to initialize database: ", err)
	} else {
		log.Println("Database initialized!")
	}
	http.HandleFunc("/", RegisterPing)
	http.HandleFunc("/pingpong", RegisterPing)
	http.HandleFunc("/healthz", Healthz)
	log.Println("Ping/pong server listening in address http://localhost" + port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
