package render

import (
	"fmt"
	"io"

	"github.com/mplewis/kesdev.com/post"
)

var postTmpl = loadTemplate("post", "templates/post.html")

type postStub struct {
	NavTitle    string
	SiteName    string
	PublishedAt string
	UpdatedAt   string
	PostTitle   string
	Body        string
}

func Post(dst io.Writer, p post.Post) error {
	stub := postStub{
		NavTitle:    fmt.Sprintf("%s – %s", p.Title, siteName),
		SiteName:    siteName,
		PublishedAt: p.CreatedAt.Format(dateFormatLong),
		UpdatedAt:   p.UpdatedAt.Format(dateFormatLong),
		PostTitle:   p.Title,
		Body:        p.HTML,
	}
	return postTmpl.Execute(dst, stub)
}
