package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"cartracker.api/common"
	"cartracker.api/data"
	"cartracker.api/routers"
	"cartracker.api/settings"
	"github.com/urfave/negroni"
)

// HTMLServer the struct
type HTMLServer struct {
	server *http.Server
	wg     sync.WaitGroup
}

// our main function
func main() {
	settings.Init()

	common.MongoSession = data.InitMongoSession()

	htmlServer := Start(common.ServerCfg)
	defer htmlServer.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("main : shutting down")
}

// Start launches the HTML Server
func Start(cfg *common.ServerConfig) *HTMLServer {
	// Setup Context
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup Handlers
	router := routers.NewRouter()

	n := negroni.Classic() // Includes some default middlewares
	n.Use(negroni.NewStatic(http.Dir("/static")))
	n.UseHandler(router)

	// Create the HTML Server
	htmlServer := HTMLServer{
		server: &http.Server{
			Addr:           ":" + cfg.Host,
			Handler:        n,
			ReadTimeout:    cfg.ReadTimeout * time.Second,
			WriteTimeout:   cfg.WriteTimeout * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	fmt.Print(cfg.ReadTimeout)

	// Add to the WaitGroup for the listener goroutine
	htmlServer.wg.Add(1)

	// Start the listener
	go func() {
		fmt.Printf("\nHTMLServer : Service started : Host=%v\n", ":"+cfg.Host)
		htmlServer.server.ListenAndServe()
		htmlServer.wg.Done()
	}()

	return &htmlServer
}

// Stop turns off the HTML Server
func (htmlServer *HTMLServer) Stop() error {
	// Create a context to attempt a graceful 5 second shutdown.
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("\nHTMLServer : Service stopping\n")

	// Attempt the graceful shutdown by closing the listener
	// and completing all inflight requests
	if err := htmlServer.server.Shutdown(ctx); err != nil {
		// Looks like we timed out on the graceful shutdown. Force close.
		if err := htmlServer.server.Close(); err != nil {
			fmt.Printf("\nHTMLServer : Service stopping : Error=%v\n", err)
			return err
		}
	}

	// Wait for the listener to report that it is closed.
	htmlServer.wg.Wait()
	fmt.Printf("\nHTMLServer : Stopped\n")
	return nil
}
