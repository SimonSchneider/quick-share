package main

import (
	"embed"
	_ "embed"
	"flag"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func Must[T any](t T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return t
}

//go:embed static/*
var embeddedFS embed.FS

func getFS(watch bool) fs.FS {
	if watch {
		if _, err := os.Stat("static"); err == nil {
			return os.DirFS("static")
		}
	}
	return Must(fs.Sub(embeddedFS, "static"))
}

type templates struct {
	files fs.FS
}

func (w *templates) Parse() (*template.Template, error) {
	return template.ParseFS(w.files, "*.gohtml")
}

func (w *templates) Lookup(pth string) *template.Template {
	tmpls, err := w.Parse()
	if err != nil {
		log.Fatal(err)
	}
	tmpl := tmpls.Lookup(pth)
	if tmpl == nil {
		log.Fatalf("template %s not found", pth)
	}
	return tmpl
}

func getTemplates(files fs.FS, watch bool) Templates {
	tmpls := &templates{files: Must(fs.Sub(files, "templates"))}
	if watch {
		return tmpls
	}
	return Must(tmpls.Parse())
}

var fWatch = flag.Bool("watch", false, "Watch for changes in the static directory (useful for debugging)")
var fAddr = flag.String("addr", ":8888", "Address to listen on")
var fMaxSecrets = flag.Uint("max-secrets", 100, "Maximum number of Secrets before evicting old ones")
var fMaxSecretSize = flag.Int64("max-secret-size", 256*1024, "Maximum size of a secret in bytes")

func main() {
	flag.Parse()
	files := getFS(*fWatch)
	tmpls := getTemplates(files, *fWatch)
	handler := Handler{
		Secrets:        NewInMemorySecrets(*fMaxSecrets),
		Templates:      tmpls,
		MaxSecretBytes: *fMaxSecretSize,
		Files:          files,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /secrets", handler.CreateSecret)
	mux.HandleFunc("GET /secrets/{id}", handler.GetSecret)
	mux.Handle("GET /", http.StripPrefix("/", http.FileServerFS(Must(fs.Sub(files, "public")))))

	log.Printf("Listening at: %s", *fAddr)
	if err := http.ListenAndServe(*fAddr, mux); err != nil {
		log.Fatal(err)
	}
}
