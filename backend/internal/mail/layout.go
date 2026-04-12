package mail

import (
	"bytes"
	"embed"
	"html/template"
	"sync"
)

// EmailLayout fills templates/layout.html: shell + main column + optional button + optional URL line.
// Body is the main column HTML (paragraphs with inline styles). Plain fields are escaped by html/template.
// No button: then leave ButtonLabel and ButtonURL empty. No URL line: then leave LinkURL empty.
type EmailLayout struct {
	AppName     string
	Badge       string
	BodyIntro   template.HTML
	BodyMuted   string
	ButtonLabel string
	ButtonURL   string
	LinkHint    string
	LinkURL     string
	Footer      string
}

//go:embed templates/*.html
var layoutFS embed.FS

var (
	layoutOnce sync.Once
	layoutTmpl *template.Template
	layoutErr  error
)

func renderEmailLayout(d EmailLayout) (string, error) {
	layoutOnce.Do(func() {
		layoutTmpl, layoutErr = template.ParseFS(layoutFS, "templates/layout.html")
	})
	if layoutErr != nil {
		return "", layoutErr
	}
	var buf bytes.Buffer
	if err := layoutTmpl.Execute(&buf, d); err != nil {
		return "", err
	}
	return buf.String(), nil
}
