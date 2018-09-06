package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"encoding/json"
	"fmt"

	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/vidmed/detection"
	"github.com/vidmed/logger"
)

var (
	configFileName   = flag.String("config", "config.toml", "Config file name")
	bidRequestParser *detection.BidRequestParser
)

func init() {
	flag.Parse()
	config, err := NewConfig(*configFileName)
	if err != nil {
		logger.Get().Fatalf("ERROR loading config: %s\n", err.Error())
	}
	// Init logging, logger goes first since other components may use it
	logger.Init(int(GetConfig().Main.LogLevel))

	parser, err := detection.NewFiftyoneDegreesParser(config.Main.FiftyOneDegreesDBPath)
	if err != nil {
		logger.Get().Fatalf("ERROR creating FiftyoneDegreesParser: %s\n", err.Error())
	}

	country_detector, err := detection.NewMaxMindCountryDetector(config.Main.MaxMinDBPath)
	if err != nil {
		logger.Get().Fatalf("ERROR creating MaxMindCountryDetector: %s\n", err.Error())
	}

	bidRequestParser = detection.NewBidRequestParser(country_detector, parser)
	logger.Get().Info("BidRequestParser inited successfully")
}

func main() {
	numCPUs := runtime.NumCPU()
	logger.Get().Infof("CPUs count %d", numCPUs)
	runServer()
}

func runServer() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	http.HandleFunc("/", recoverHandler(handler))
	hs := &http.Server{Addr: fmt.Sprintf("%s:%d", GetConfig().Main.ListenAddr, GetConfig().Main.ListenPort), Handler: nil}

	go func() {
		logger.Get().Infof("Listening on http://%s\n", hs.Addr)

		if err := hs.ListenAndServe(); err != http.ErrServerClosed {
			logger.Get().Fatal(err.Error())
		}
	}()

	<-stop

	timeout := 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	fmt.Printf("Shutdown with timeout: %s\n", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		logger.Get().Errorf("Error: %v\n", err)
	} else {
		logger.Get().Infof("Server stopped")
	}
	cancel()
}

func handler(w http.ResponseWriter, r *http.Request) {
	// base checks
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	var request = new(detection.BidRequest)
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		logger.Get().Errorln(err)
		http.Error(w, "error decoding request body", http.StatusInternalServerError)
		return
	}

	response, err := bidRequestParser.Parse(request)
	if err != nil {
		logger.Get().Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func recoverHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if errR := recover(); errR != nil {
				err := errors.Errorf("panic: %+v", errR)
				logger.Get().Error(err)
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		handler(w, r)
	}
}
