package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

type Secrets interface {
	Add(encryptedSecret string, uses int, exp time.Duration) string
	Get(id string) (string, bool)
}

type Templates interface {
	Lookup(name string) *template.Template
}

type Handler struct {
	Secrets        Secrets
	Templates      Templates
	Files          fs.FS
	MaxSecretBytes int64
}

type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}

type SecretRequest struct {
	EncryptedSecret string   `json:"encrypted_secret"`
	Uses            int      `json:"uses"`
	Expiration      Duration `json:"expiration"`
}

func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	var req SecretRequest
	limited := http.MaxBytesReader(w, r.Body, h.MaxSecretBytes)
	defer limited.Close()
	if err := json.NewDecoder(limited).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := h.Secrets.Add(req.EncryptedSecret, req.Uses, time.Duration(req.Expiration))

	w.Write([]byte(id))
}

func (h *Handler) execHtmlTemplate(w http.ResponseWriter, name string, data any) {
	tmpl := h.Templates.Lookup(name)
	if tmpl == nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	respType := r.Header.Get("Content-Type")

	encryptedSecret, ok := h.Secrets.Get(id)
	if !ok {
		if respType == "text/html" || respType == "" {
			http.ServeFileFS(w, r, h.Files, "secretNotFound.html")
		} else {
			http.Error(w, "Secret not found", http.StatusNotFound)
		}
		return
	}
	if respType == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"encrypted_secret": encryptedSecret})
	} else if respType == "text/html" || respType == "" {
		h.execHtmlTemplate(w, "secrets.gohtml", struct {
			EncryptedSecret string
		}{
			EncryptedSecret: encryptedSecret,
		})
	} else if respType == "text/plain" {
		w.Write([]byte(encryptedSecret))
	} else {
		http.Error(w, fmt.Sprintf("Unsupported content type: '%s'", respType), http.StatusUnsupportedMediaType)
	}
}
