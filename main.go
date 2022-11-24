package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	// os.Exit関数が呼ばれないようにcontext.Contextを受け取る
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context) error {
	// *http.Server型を使うとサーバーのタイムアウト時間など柔軟に設定できる
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	}
	// errgruup.Group型を使うと戻り値にエラーが含まれるゴルーチン間の並行処理の実装が簡単になる
	eg, ctx := errgroup.WithContext(ctx)

	// 別ゴルーチンでHTTPサーバーを起動
	eg.Go(func() error {
		// http.ErrServerClosedはhttp.Server.Shatdown()が正常に終了したことを示す
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
