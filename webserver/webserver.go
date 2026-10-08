package webserver

import (
	"audiobook-ingest/config"
	"audiobook-ingest/processor"
	"audiobook-ingest/webserver/websocket"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	gowebsocket "golang.org/x/net/websocket"
)

var (
	htmlTemplate    *template.Template
	embedFilesystem embed.FS
	rawFaviconSVG   []byte
)

func InitTemplates(fs embed.FS) {
	var err error

	funcMap := template.FuncMap{
		"json": func(v any) string {
			bytes, err := json.Marshal(v)
			if err != nil {
				return "[]"
			}
			return string(bytes)
		},
	}

	htmlTemplate, err = template.New("").Funcs(funcMap).ParseFS(fs,
		"html/index.html",
		"html/partials/*.html",
		"html/partials/sections/*.html",
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to load dashboard and partial templates: %v", err))
	}

	rawFaviconSVG, err = fs.ReadFile("html/favicon.svg")
	if err != nil {
		panic(fmt.Sprintf("Failed to map favicon asset from embed array: %v", err))
	}

	embedFilesystem = fs
	websocket.Setup()
}

func subFS(embeddedFS embed.FS, subDir string) http.FileSystem {
	sub, err := fs.Sub(embeddedFS, subDir)
	if err != nil {
		panic(fmt.Sprintf("Failed to extract sub-filesystem path link for %s: %v", subDir, err))
	}
	return http.FS(sub)
}

func Start() {
	http.HandleFunc("/", handleDashboard)

	http.Handle("/css/", http.StripPrefix("/css/", staticMimeMux(http.FileServer(subFS(embedFilesystem, "html/css")))))
	http.Handle("/js/", http.StripPrefix("/js/", staticMimeMux(http.FileServer(subFS(embedFilesystem, "html/js")))))
	http.Handle("/img/", http.StripPrefix("/img/", staticMimeMux(http.FileServer(subFS(embedFilesystem, "html/img")))))

	http.Handle("/ws", gowebsocket.Handler(handleWebsocket))

	// IDEA can this be made to be an template so we can change the text and possibly make it say [ext] to [ext]
	http.HandleFunc("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write(rawFaviconSVG)
	})

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
	}

	websocket.Start()
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templateContext := map[string]any{
		"ShowWarningBanner":  !config.HideAPIWarning && !processor.Process.APISearchEnabled(),
		"RequiredEnvKeys":    processor.Process.EnvKeyNames(),
		"Instructions":       processor.Process.EnvInstructions(),
		"ManualSearchFields": processor.Process.GetManualSearchFields(),
		"DisableAutoIngest":  config.DisableAutoIngest,
		"MultiSwitchCount":   config.MultiSwitchCount,
	}

	if err := htmlTemplate.ExecuteTemplate(w, "index.html", templateContext); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// IDEA Look to see if this can be removed later on
func staticMimeMux(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(r.URL.Path)
		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case ".svg":
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

func handleWebsocket(ws *gowebsocket.Conn) {
	websocket.AddWebsocket(ws)
	websocket.HandleWebsocketStream(ws)
}
