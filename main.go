package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/k0kubun/pp/v3"
	"github.com/mplewis/kesdev.com/post"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func verifyNoDuplicates(posts []post.Post) error {
	slugCount := map[string]int{}
	for _, p := range posts {
		slugCount[p.Slug]++
	}
	dupes := []string{}
	for slug, count := range slugCount {
		if count > 1 {
			dupes = append(dupes, slug)
		}
	}
	if len(dupes) > 0 {
		return fmt.Errorf("duplicate slugs: %s", strings.Join(dupes, ", "))
	}
	return nil
}

func main() {
	files, err := filepath.Glob("content/*.md")
	check(err)
	posts := make([]post.Post, len(files))
	for i, file := range files {
		f, err := os.Open(file)
		check(err)
		defer f.Close()
		p, err := post.Parse(f)
		check(err)
		posts[i] = p
	}

	check(verifyNoDuplicates(posts))

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt) // newest first
	})

	pp.Println(posts)
}
