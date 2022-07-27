package render

import (
	"io"
	"os"
	"strings"
	"text/template"
)

// TODO: Configize
const siteName = "Kestrel Development"
const postPathPrefix = "/post"
const yieldContentSigil = "<!-- yield content -->"
const dateFormatShort = "2006-01-02"
const dateFormatLong = "2006-01-02 15:04 MST"

var layout = load("templates/layout.html")

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(path string) string {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	raw, err := io.ReadAll(f)
	check(err)
	return string(raw)
}

func loadTemplate(name string, path string) *template.Template {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	raw, err := io.ReadAll(f)
	check(err)
	content := string(raw)
	complete := strings.ReplaceAll(layout, yieldContentSigil, content)
	return template.Must(template.New(name).Parse(complete))
}
