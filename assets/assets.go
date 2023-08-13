package assets

import (
	"bytes"
	"embed"
	"github.com/sharat87/httpbun/exchange"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//go:embed *.html *.css *.js *.png *.svg favicon.ico site.webmanifest
var assets embed.FS

func Render(name string, ex exchange.Exchange, data map[string]any) {
	data["serverSpec"] = ex.ServerSpec

	ex.ResponseWriter.Header().Set("Content-Type", "text/html")

	tpl, err := template.ParseFS(assets, "*")
	if err != nil {
		log.Fatalf("Error parsing HTML assets %v.", err)
	}

	var rendered bytes.Buffer
	if err := tpl.ExecuteTemplate(&rendered, name, data); err != nil {
		log.Fatalf("Error executing %q template %v.", name, err)
	}

	_, err = ex.ResponseWriter.Write(rendered.Bytes())
	if err != nil {
		log.Printf("Error writing rendered template %v", err)
	}
}

func WriteAsset(name string, w http.ResponseWriter, req *http.Request) {
	if content, err := assets.ReadFile(name); err == nil {
		_, err := w.Write(content)
		if err != nil {
			log.Printf("Error writing asset content %v", err)
		}
	} else if strings.HasSuffix(err.Error(), " file does not exist") {
		http.NotFound(w, req)
	} else {
		log.Printf("Error %v", err)
	}
}
