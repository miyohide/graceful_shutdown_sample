# Graceful Shutdownのサンプル

## 概要

Graceful Shutdownをどのように使うのかがよく分からなかったので、サンプルを書いてみました。

https://qiita.com/najeira/items/806cacb9bba96ff06ec4 を参考にしました。

## 使い方

以下のコマンドを実行します。

- go build -o gs.exe
- go get github.com/lestrrat-go/server-starter/cmd/start_server
- $GOPATH/bin/start_server --port 8080 --pid-file app.pid -- ./gs.exe

HTMLリクエストは `curl http://127.0.0.1:8080/` で送ります。

サーバの再起動は

```
kill -HUP `cat app.pid`
```

サーバの停止は

```
kill -TERM `cat app.pid`
```

Graceful Shutdownを確認するには、例えば、サーバを起動させた状態で

```
while true; do kill -HUP `cat app.pid`; sleep 1; done
```

というふうに再起動処理を1秒間隔で実施しているなかでabなんかでアクセスしても処理が失敗しないことが確認できる
