package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/{{githubAccount}}/{{name}}/internal/routes"
	"github.com/{{githubAccount}}/{{name}}/pkg/logger"
)

type PrettifyWriter struct {
	writer *os.File
}

func (p *PrettifyWriter) Write(data []byte) (int, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return p.writer.Write(data)
	}

	prettyJSON.WriteByte('\n')
	return p.writer.Write(prettyJSON.Bytes())
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := logger.GetJsonLogger()

	srv := newServer(logger)
	httpServer := &http.Server{
		Addr:     ":8080",
		Handler:  srv,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	go func() {
		logger.Info(fmt.Sprintf("listening on %s\n", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Error(fmt.Sprintf("error listening and serving: %s", err))
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownContext := context.Background()
		shutdownContext, cancel := context.WithTimeout(shutdownContext, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownContext); err != nil {
			logger.Error(fmt.Sprintf("error shutting down http server: %s", err))
		}
	}()

	wg.Wait()
	return nil
}

func newServer(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(mux, logger)

	var handler http.Handler = mux
	return handler
}
