package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type Site struct {
	Path     string
	Output   string
	Template *template.Template
	Rules    map[string]([]string)
	Pages    PageSlice
}

func NewSite(config *GlobalConfig) *Site {
	template, err := template.ParseFiles(config.Templates...)
	errhandle(err)

	site := &Site{
		Path: config.Source,
		Output: config.Output,
		Template: template,
		Rules: config.Rules,
		Pages: make(PageSlice, 0),
	}

	site.Collect()

	return site
}

func (site *Site) AddPage(path string) {
	page := NewPage(site, path)
	site.Pages = append(site.Pages, page)
}

func (site *Site) Collect() {
	errors := make(chan error)

	filepath.Walk(site.Path, site.walkFunc(errors))

	select {
	case err := <-errors:
		errhandle(err)
	default:
	}
}

func (site *Site) walkFunc(errors chan<- error) filepath.WalkFunc {
	return func(fn string, fi os.FileInfo, err error) error {
		if err != nil {
			errors <- err
			return nil
		}

		if !fi.IsDir() {
			site.AddPage(fn)
		}

		return nil
	}
}

func (site *Site) Summary() {
	println("Total pages", len(site.Pages))
	for _, page := range site.Pages {
		fmt.Printf("%s - %s: %d chars\n",
			page.Path, page.Title, len(page.Content))
		println("------------")
		fmt.Printf("%s", page.Rendered())
		println("------------")
		println()
	}
}