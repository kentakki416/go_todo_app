package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"github.com/kentakki416/go_todo_app/handler"
	"github.com/kentakki416/go_todo_app/store"
)

func NewMux() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		//静的解析のエラーを回避するために明示的に戻り値を捨てている
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	v := validator.New()
	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post("/tasks:", at.ServeHTTP)
	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get("/tasks", lt.ServeHTTP)
	return mux
}
