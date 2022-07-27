package post

import (
	"bytes"
	"errors"
	"io"
	"time"

	fmlib "github.com/adrg/frontmatter"
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

func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("missing title")
	}
	if p.HTML == "" {
		return errors.New("missing html")
	}
	return nil
}

func (p *Post) GenerateSlug() string {
	return p.Title // TODO
}

func (p *Post) GenerateExcerpt() string {
	return p.HTML // TODO
}

func (p *Post) Fill() error {
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = p.CreatedAt
	}
	if p.PublishedAt.IsZero() {
		p.PublishedAt = p.CreatedAt
	}
	if p.Excerpt == "" {
		p.Excerpt = p.GenerateExcerpt()
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
