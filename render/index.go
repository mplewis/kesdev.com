package render

import (
	"fmt"
	"io"

	"github.com/mplewis/kesdev.com/post"
)

var indexTmpl = loadTemplate("index", "templates/index.html")

type indexStub struct {
	Date      string
	Link      string
	PostTitle string
	Excerpt   string
}

type indexArgs struct {
	NavTitle string
	SiteName string
	Stubs    []indexStub
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
