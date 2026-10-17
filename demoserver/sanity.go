// Sanity demo for coreserver. Run with `go run .` from demoserver, then try:
//
//	curl localhost:8080/health
//	curl "localhost:8080/hello?name=gopher"
//	curl localhost:8080/panic                         # 500, panic logged with stack
//	curl -X POST "localhost:8080/loglevel?level=debug" # /panic now returns the panic message
//	curl localhost:8080/slow & sleep 0.5; kill -INT <pid> # request finishes before exit
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/isingh6022/go-fullstack-ai/coreserver/server"
)

func main() {
	app := server.CreateApp(8080, server.LevelInfo)

	app.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	app.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "world"
		}
		app.Logger.Debug("hello called", "name", name)
		fmt.Fprintf(w, "hello, %s\n", name)
	})

	app.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
		panic("demo panic")
	})

	app.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info(fmt.Sprintf("slow request started at: %v", time.Now().Second()))
		time.Sleep(3 * time.Second)
		
		w.Write([]byte("slow done\n"))
		app.Logger.Info(fmt.Sprintf("slow request completed at: %v", time.Now().Second()))
	})

	app.HandleFunc("POST /loglevel", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("level") {
		case "debug":
			app.Logger.SetLevel(server.LevelDebug)
		case "info":
			app.Logger.SetLevel(server.LevelInfo)
		case "warn":
			app.Logger.SetLevel(server.LevelWarn)
		case "error":
			app.Logger.SetLevel(server.LevelError)
		default:
			http.Error(w, "level must be one of: debug, info, warn, error", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "log level set to %s\n", app.Logger.Level())
	})

	app.HandleFunc("POST /shutdown", func(w http.ResponseWriter, r*http.Request){
		w.Write([]byte("shutting down\n"))
		go app.Shutdown()
	})

	if err := app.Start(); err != nil {
		app.Logger.Error("server exited", "error", err)
		os.Exit(1)
	}
}
