package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/SimonSchneider/goslu/config"
	"github.com/SimonSchneider/goslu/srvu"
	"github.com/SimonSchneider/goslu/templ"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

//go:embed static/*
var embeddedFS embed.FS

type Config struct {
	Watch         bool
	Addr          string
	MaxSecrets    uint
	MaxSecretSize int64 `config:"u: Maximum size of a secret in bytes"`
}

func parseConfig(args []string, getEnv func(string) string) (Config, error) {
	cfg := Config{
		Addr:          ":8888",
		MaxSecrets:    100,
		MaxSecretSize: 256 * 1024,
	}
	return cfg, config.ParseInto(&cfg, flag.NewFlagSet("", flag.ExitOnError), args, getEnv)
}

func Run(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer, getEnv func(string) string, getwd func() (string, error)) error {
	cfg, err := parseConfig(args[1:], getEnv)
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	logger := srvu.LogToOutput(log.New(stdout, "", log.LstdFlags|log.Lshortfile))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	pub, tmpl, err := templ.GetPublicAndTemplates(embeddedFS, &templ.Config{
		Watch:        cfg.Watch,
		TmplPatterns: []string{"templates/*.gohtml"},
	})
	if err != nil {
		return fmt.Errorf("failed to get public and templates: %w", err)
	}
	fmt.Printf("pub: %v, tmpl: %v\n", pub, tmpl.Lookup("secrets.gohtml"))
	handler := Handler{
		Secrets:        NewInMemorySecrets(cfg.MaxSecrets),
		Templates:      tmpl,
		MaxSecretBytes: cfg.MaxSecretSize,
		Files:          pub,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /secrets", handler.CreateSecret)
	mux.HandleFunc("GET /secrets/{id}", handler.GetSecret)
	mux.Handle("GET /", http.StripPrefix("/", http.FileServerFS(pub)))

	srv := &http.Server{
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		Addr:    cfg.Addr,
		Handler: srvu.With(mux, srvu.WithCompression(), srvu.WithLogger(logger)),
	}
	logger.Printf("Listening on: %s", cfg.Addr)
	return srvu.RunServerGracefully(ctx, srv, logger)
}
