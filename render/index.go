package render

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/mplewis/kesdev.com/post"
)

const siteName = "Kestrel Development"
const postPathPrefix = "/post"
const yieldContentSigil = "<!-- yield content -->"
const dateFormatShort = "2006-01-02"
const dateFormatLong = "2006-01-02 15:04 MST"

var layout = load("templates/layout.html")
var indexTmpl = loadTemplate("index", "templates/index.html")
var postTmpl = loadTemplate("post", "templates/post.html")

type indexStub struct {
	Date      string
	Link      string
	PostTitle string
	Excerpt   string
}

type postStub struct {
	NavTitle    string
	SiteName    string
	PublishedAt string
	UpdatedAt   string
	PostTitle   string
	Body        string
}

type indexArgs struct {
	NavTitle string
	SiteName string
	Stubs    []indexStub
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
	stubs := make([]indexStub, len(posts))
	for i, p := range posts {
		stubs[i] = indexStub{
			Date:      p.CreatedAt.Format(dateFormatShort),
			Link:      fmt.Sprintf("%s/%s", postPathPrefix, p.Slug),
			PostTitle: p.Title,
			Excerpt:   p.Excerpt,
		}
	}
	args := indexArgs{NavTitle: siteName, SiteName: siteName, Stubs: stubs}
	return indexTmpl.Execute(dst, args)
}

func Post(dst io.Writer, p post.Post) error {
	const emDash = "—" // make linter happy
	stub := postStub{
		NavTitle:    fmt.Sprintf("%s %s %s", p.Title, emDash, siteName),
		SiteName:    siteName,
		PublishedAt: p.CreatedAt.Format(dateFormatLong),
		UpdatedAt:   p.UpdatedAt.Format(dateFormatLong),
		PostTitle:   p.Title,
		Body:        p.HTML,
	}
	return postTmpl.Execute(dst, stub)
}
