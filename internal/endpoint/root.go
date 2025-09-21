package endpoint

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/Dydhzo/stremthru-proxy/internal/config"
	"github.com/Dydhzo/stremthru-proxy/internal/shared"
)

//go:embed root.html
var templateBlob string

type rootTemplateDataSection struct {
	Title   string        `json:"title"`
	Content template.HTML `json:"content"`
}

type RootTemplateData struct {
	Title       string                    `json:"-"`
	Description template.HTML             `json:"description"`
	Version     string                    `json:"-"`
	Sections    []rootTemplateDataSection `json:"sections"`
}

var rootTemplateData = func() RootTemplateData {
	td := RootTemplateData{}
	err := json.Unmarshal([]byte(config.LandingPage), &td)
	if err != nil {
		panic("malformed config for landing page: " + config.LandingPage)
	}
	return td
}()

var ExecuteTemplate = func() func(data *RootTemplateData) (bytes.Buffer, error) {
	tmpl := template.Must(template.New("root.html").Parse(templateBlob))
	return func(data *RootTemplateData) (bytes.Buffer, error) {
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, data)
		return buf, err
	}
}()

func handleRoot(w http.ResponseWriter, r *http.Request) {
	td := &RootTemplateData{
		Title:       "StremThru Proxy",
		Description: rootTemplateData.Description,
		Version:     config.Version,
		Sections:    rootTemplateData.Sections,
	}

	buf, err := ExecuteTemplate(td)
	if err != nil {
		shared.SendError(w, r, err)
		return
	}
	shared.SendHTML(w, 200, buf)
}

func AddRootEndpoint(mux *http.ServeMux) {
	mux.HandleFunc("/{$}", handleRoot)
}
