package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"go-url-shortener"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	err := godotenv.Load()
	if err != nil {
		logger.Log("Could not load environment file.")
	}

	var (
		httpAddr = flag.String("http.addr", os.Getenv("APP_PORT"), "HTTP listen address")
	)
	flag.Parse()

	var s shortener.Service
	{
		s = shortener.NewInmemService()
		s = shortener.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = shortener.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
