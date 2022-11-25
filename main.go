package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}
	// os.Exit関数が呼ばれないようにcontext.Contextを受け取る
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error {
	// *http.Server型を使うとサーバーのタイムアウト時間など柔軟に設定できる
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	}
	// errgruup.Group型を使うと戻り値にエラーが含まれるゴルーチン間の並行処理の実装が簡単になる
	eg, ctx := errgroup.WithContext(ctx)

	// 別ゴルーチンでHTTPサーバーを起動
	eg.Go(func() error {
		// http.ErrServerClosedはhttp.Server.Shatdown()が正常に終了したことを示す
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの通知（終了通知）を待機する
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdwn: %+v", err)
	}

	// Goメソッドで起動した別のゴルーチンの終了を待つ
	return eg.Wait()
}
