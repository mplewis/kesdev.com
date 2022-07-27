package post

import (
	"errors"
	"time"
)

type Post struct {
	Title       string
	Slug        string // optional
	Published   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time // optional
	PublishedAt time.Time // optional
	HTML        string
	Excerpt     string // optional
}

func (p *Post) Validate() error {
	if p.Title != "" {
		return errors.New("missing title")
	}
	if p.HTML != "" {
		return errors.New("missing html")
	}
	return nil
}

func (p *Post) GenerateSlug() string {
	return p.Title
}

func (p *Post) GenerateExcerpt() string {
	return p.HTML // TODO
}

func (p *Post) Fill() error {
	if err := p.Validate(); err != nil {
		return err
	}
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

func Parse(raw string) (Post, error) {
	return Post{}, nil
}
