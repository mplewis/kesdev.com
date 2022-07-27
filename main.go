package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mplewis/kesdev.com/post"
	"github.com/mplewis/kesdev.com/render"
	"github.com/yargevad/filepathx"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func openWrite(path string) io.Writer {
	os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	check(err)
	return f
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

func clean() {
	toDel, err := filepathx.Glob("dist/**/*.html")
	check(err)
	for _, f := range toDel {
		check(os.Remove(f))
	}
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

	clean()

	f := openWrite(path.Join("dist", "index.html"))
	check(render.Index(f, posts))
	fmt.Println("Wrote dist/index.html")

	for _, p := range posts {
		// TODO: Configize posts path
		dest := path.Join("dist", "post", p.Slug, "index.html")
		f := openWrite(dest)
		render.Post(f, p)
		fmt.Printf("Wrote %s\n", dest)
	}
}
