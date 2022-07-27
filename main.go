package main

import (
	"os"
	"path/filepath"

	"github.com/k0kubun/pp/v3"
	"github.com/mplewis/kesdev.com/post"
)

func check(err error) {
	if err != nil {
		panic(err)
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
	pp.Println(posts)
}
