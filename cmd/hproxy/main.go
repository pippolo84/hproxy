package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Pippolo84/hproxy/internal/hedged"
	"github.com/Pippolo84/hproxy/internal/service"
	"github.com/gorilla/mux"
)

const (
	// Name is the name of the service
	Name string = "hrpoxy"

	// Address is the string representation of the address where the service will listen
	Address string = ":8081"

	// SvcCooldownTimeout is the maximum cooldown time before forcing the shutdown
	SvcCooldownTimeout time.Duration = 10 * time.Second
)

func main() {
	proxy := NewProxy()

	if err := proxy.Init(); err != nil {
		log.Fatalf("initialization error: %v\n", err)
	}
	defer proxy.Close()

	errs := proxy.Run()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// block until a signal or an error is received
	select {
	case err := <-errs:
		log.Println(err)
	case sig := <-signalChan:
		log.Printf("got signal: %v, shutting down...\n", sig)
	}

	// graceful shutdown the service
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), SvcCooldownTimeout)
	defer cancelShutdown()

	var stopWg sync.WaitGroup
	stopWg.Add(1)

	if err := proxy.Shutdown(shutdownCtx, &stopWg); err != nil {
		log.Println(err)
	}

	// wait for service to cleanup everything
	stopWg.Wait()
}

// Proxy represents our specific service
type Proxy struct {
	*service.Service
	client *hedged.Client
}

// NewProxy returns a new Proxy ready to be initialized and run
func NewProxy() *Proxy {
	proxy := &Proxy{}

	router := mux.NewRouter()

	router.HandleFunc("/", proxy.handler).Schemes("http").Methods(http.MethodGet)

	proxy.Service = service.NewDefaultService(Name, Address, router)

	proxy.client = hedged.NewClient(http.DefaultClient, 80*time.Second, 3)

	return proxy
}

// Init initializes private components of the service, like its internal proxy
func (p *Proxy) Init() error {
	rand.Seed(time.Now().UnixNano())

	return nil
}

// Close frees all resources used by the service private components
func (p *Proxy) Close() {
}

// ************************** HTTP Handlers **************************

func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/", nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, err := p.client.Do(req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
