// Package grace
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 22:00
package grace

import (
	"log"
	"net/http"
	"syscall"
	"testing"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WORLD!"))
}

func TestNewServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handler)

	srv := NewServer(":8080", mux)

	srv.RegisterSignalHook(PreSignal, syscall.SIGINT, func() {
		t.Log("之前")
	})
	srv.RegisterSignalHook(PostSignal, syscall.SIGINT, func() {
		t.Log("之后")
	})

	//err := ListenAndServe("localhost:8080", mux)
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
	log.Println("Server on 8080 stopped")
	//os.Exit(0)
}
