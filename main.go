package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/lestrrat-go/server-starter/listener"
)

func allPathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, Graceful Shutdown Sample")
	w.(http.Flusher).Flush()
}

func newHandler() http.Handler {
	mux := http.NewServeMux()
	conf, err := config()
	if err != nil {
		log.Fatalln("Raise error when get data from redis")
	}
	fmt.Printf("Starting new Serve. config val = %v\n", conf)
	mux.HandleFunc("/", allPathHandler)
	return mux
}

func config() (result string, err error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return "", err
	}
	defer c.Close()

	result, err = redis.String(c.Do("GET", "CONFIG"))
	if err != nil {
		return "", err
	}
	return result, err
}

func main() {
	listeners, err := listener.ListenAll()
	if err != nil && err != listener.ErrNoListeningTarget {
		log.Fatal(err)
	}
	// http.Server構造体の初期化
	server := &http.Server{Handler: newHandler()}

	// サーバはブロックするので別のgoroutineで実行する
	go func() {
		if err := server.Serve(listeners[0]); err != nil {
			log.Print(err)
		}
	}()

	// シグナルを待つ
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

	// シグナルを受け取ったらShutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Print(err)
	}
}
