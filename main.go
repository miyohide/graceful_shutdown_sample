// Graceful Shutdownのサンプル
// https://qiita.com/najeira/items/806cacb9bba96ff06ec4 を参考にした
// 使い方
// - go build -o gs.exe
// - go get github.com/lestrrat-go/server-starter/cmd/start_server
// - $GOPATH/bin/start_server --port 8080 --pid-file app.pid -- ./gs.exe
// HTMLリクエストは curl http://127.0.0.1:8080/ で送る
// サーバの再起動は kill -HUP `cat app.pid`
// サーバの停止は kill -TERM `cat app.pid`
// 例えば、サーバを起動させた状態で
// while true; do kill -HUP `cat app.pid`; sleep 1; done
// というふうに再起動処理を1秒間隔で実施しているなかで
// abなんかでアクセスしても処理が失敗しないことが確認できる
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

	"github.com/lestrrat-go/server-starter/listener"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, ")
	w.(http.Flusher).Flush()
	time.Sleep(time.Millisecond * 100)
	fmt.Fprint(w, "Go 1.8!\n")
}

func main() {
	listeners, err := listener.ListenAll()
	if err != nil && err != listener.ErrNoListeningTarget {
		log.Fatal(err)
	}
	// http.Server構造体の初期化
	server := &http.Server{Handler: http.HandlerFunc(handler)}

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
