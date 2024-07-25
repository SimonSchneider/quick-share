package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func NewTemplateProvider(watch bool) SecretTemplateProvider {
	pth := filepath.Join("templates", "secret.html")
	if watch {
		return func() (*template.Template, error) {
			return template.ParseFiles(pth)
		}
	}
	tmpl := template.Must(template.ParseFiles(pth))
	return func() (*template.Template, error) {
		return tmpl, nil
	}
}

var fWatch = flag.Bool("watch", false, "Watch for changes in the static directory (useful for debugging)")
var fAddr = flag.String("addr", ":8888", "Address to listen on")
var fMaxSecrets = flag.Uint("max-secrets", 100, "Maximum number of Secrets before evicting old ones")
var fMaxSecretSize = flag.Int64("max-secret-size", 256*1024, "Maximum size of a secret in bytes")

func main() {
	flag.Parse()
	handler := Handler{
		Secrets:        NewInMemorySecrets(*fMaxSecrets),
		SecretTemplate: NewTemplateProvider(*fWatch),
		MaxSecretBytes: *fMaxSecretSize,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /create", handler.CreateSecret)
	mux.HandleFunc("GET /secret/{id}", handler.GetSecret)
	mux.Handle("GET /", http.FileServer(http.Dir("./static")))

	log.Printf("Listening at: %s", *fAddr)
	if err := http.ListenAndServe(*fAddr, mux); err != nil {
		log.Fatal(err)
	}
}
