package render

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/mplewis/kesdev.com/post"
)

const postPathPrefix = "/post"
const yieldContentSigil = "<!-- yield content -->"

var layout = load("templates/layout.html")
var indexTmpl = loadTemplate("index", "templates/index.html")

type IndexStub struct {
	Date  string
	Link  string
	Title string
}

type IndexArgs struct {
	Title string
	Stubs []IndexStub
}

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

func Index(dst io.Writer, posts []post.Post) error {
	stubs := make([]IndexStub, len(posts))
	for i, p := range posts {
		stubs[i] = IndexStub{
			Date:  p.CreatedAt.Format("2006-01-02"),
			Link:  fmt.Sprintf("%s/%s", postPathPrefix, p.Slug),
			Title: p.Title,
		}
	}

	args := IndexArgs{
		Title: "Kestrel Development",
		Stubs: stubs,
	}
	return indexTmpl.Execute(dst, args)
}
