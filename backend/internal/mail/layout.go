package mail

import (
	"bytes"
	"embed"
	"html"
	"html/template"
	"strings"
	"sync"
)

// EmailLayout fills templates/layout.html: shell + main column + optional button + optional URL line.
// Body is the main column HTML (paragraphs with inline styles). Plain fields are escaped by html/template.
// No button: leave ButtonLabel and ButtonURL empty. No URL line: leave LinkURL empty.
type EmailLayout struct {
	AppName     string
	Badge       string
	Body        template.HTML
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

// MailBody is a short helper for the usual two paragraphs (HTML intro + plain muted line).
func MailBody(htmlLine template.HTML, mutedPlain string) template.HTML {
	var b strings.Builder
	b.WriteString(`<p style="margin:0 0 16px;font-size:15.2px;color:#111111;">`)
	b.WriteString(string(htmlLine))
	b.WriteString(`</p><p style="margin:0 0 20px;font-size:14px;color:#5c5c5c;">`)
	b.WriteString(html.EscapeString(mutedPlain))
	b.WriteString(`</p>`)
	return template.HTML(b.String())
}
