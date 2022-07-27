package post

import (
	"bytes"
	"errors"
	"io"
	"regexp"
	"strings"
	"time"

	fmlib "github.com/adrg/frontmatter"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

type Post struct {
	Title       string    `yaml:"Title"`
	Slug        string    `yaml:"Slug"`      // optional
	Published   bool      `yaml:"Published"` // optional
	CreatedAt   time.Time `yaml:"CreatedAt"`
	UpdatedAt   time.Time `yaml:"UpdatedAt"`   // optional
	PublishedAt time.Time `yaml:"PublishedAt"` // optional
	HTML        string
	Excerpt     string
}

const excerptMark = "✂️"
const excerptLength = 200 // characters
const slugJoiner = "-"

var wsMatcher = regexp.MustCompile(`\s+`)
var nonSlugMatcher = regexp.MustCompile(`[^a-z0-9]+`)

func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("missing title")
	}
	if p.HTML == "" {
		return errors.New("missing html")
	}
	return nil
}

func sanitizeSlug(title string) string {
	title = strings.ToLower(title)
	title = nonSlugMatcher.ReplaceAllString(title, slugJoiner)
	title = strings.Trim(title, slugJoiner)
	return title
}

func generateExcerpt(html string) string {
	pcy := bluemonday.StrictPolicy()
	text := pcy.Sanitize(html)
	text = wsMatcher.ReplaceAllString(text, " ")

	if strings.Contains(text, excerptMark) {
		return strings.TrimSpace(strings.Split(text, excerptMark)[0])
	}

	charCount := 0
	excerpt := []string{}
	all := true
	for _, word := range strings.Split(text, " ") {
		if charCount+len(word) > excerptLength {
			all = false
			break
		}
		excerpt = append(excerpt, word)
		charCount += len(word) + 1
	}

	body := strings.Join(excerpt, " ")
	if !all {
		body += "…"
	}
	return body
}

func (p *Post) Fill() error {
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = p.CreatedAt
	}
	if p.PublishedAt.IsZero() {
		p.PublishedAt = p.CreatedAt
	}
	if p.Slug == "" {
		p.Slug = sanitizeSlug(p.Title)
	} else {
		p.Slug = sanitizeSlug(p.Slug)
	}
	if p.Excerpt == "" {
		p.Excerpt = generateExcerpt(p.HTML)
	}
	return nil
}

func Parse(raw io.Reader) (Post, error) {
	p := Post{}
	md, err := fmlib.Parse(raw, &p)
	if err != nil {
		return p, err
	}

	buf := bytes.Buffer{}
	if err := goldmark.Convert(md, &buf); err != nil {
		return p, err
	}
	p.HTML = buf.String()

	if err := p.Fill(); err != nil {
		return p, err
	}

	if err := p.Validate(); err != nil {
		return p, err
	}

	return p, nil
}
