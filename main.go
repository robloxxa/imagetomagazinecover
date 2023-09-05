package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var (
	MAX_WORKERS           = 0
	PUBLIC_URL            = ""
	ADDR                  = "localhost"
	PORT                  = "8005"
	VIEWPORT_WIDTH  int64 = 848
	VIEWPORT_HEIGHT int64 = 1200
)

var ScreenshotJobQueue chan ScreenshotJob

func init() {
	var (
		ok  bool
		err error
	)
	PORT, ok = os.LookupEnv("PORT")
	if !ok {
		PORT = "3000"
	}
	PUBLIC_URL, ok = os.LookupEnv("PUBLIC_URL")
	if !ok {
		PUBLIC_URL = "http://mary.robloxxa.ru:"
	}

	workers, ok := os.LookupEnv("MAX_WORKERS")
	if ok {
		MAX_WORKERS, err = strconv.Atoi(workers)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ScreenshotJobQueue = make(chan ScreenshotJob, MAX_WORKERS*2)
	dispatcher := NewScreenshotDispatcher(MAX_WORKERS)
	chromeCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dispatcher.Run(chromeCtx)

	server := &http.Server{Addr: ":" + PORT, Handler: setupRoutes()}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	log.Printf("client ready!\nPort: %s\nWorkers: %d\nPublic url: %s", PORT, MAX_WORKERS, PUBLIC_URL)
	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	// Wait for server context to be stopped
	<-serverCtx.Done()
}
