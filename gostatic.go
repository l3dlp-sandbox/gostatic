// (c) 2012 Alexander Solovyov
// under terms of ISC license

package main

import (
	"fmt"
	"path/filepath"
	// "text/template"
	"io/ioutil"
	"encoding/json"
	goopt "github.com/droundy/goopt"
)

var Version = "0.1"

var Summary = `gostatic path/to/config.json

Build a site.
`

var showVersion = goopt.Flag([]string{"-v", "--version"}, []string{},
	"show version and exit", "")

type GlobalConfig struct {
	Templates []string
	Source string
	Output string
	Rules map[string]([]string)
}

func RetrieveGlobalConfig(path string) *GlobalConfig {
	conftext, err := ioutil.ReadFile(path)
	errhandle(err)

	var config GlobalConfig
	err = json.Unmarshal(conftext, &config)
	errhandle(err)

	basepath, _ := filepath.Split(path)
	config.Source = filepath.Join(basepath, config.Source)
	config.Output = filepath.Join(basepath, config.Output)

	templates := make([]string, len(config.Templates))
	for i, template := range config.Templates {
		templates[i] = filepath.Join(basepath, template)
	}
	config.Templates = templates

	return &config
}

func main() {
	goopt.Version = Version
	goopt.Summary = Summary

	goopt.Parse(nil)

	if *showVersion {
		fmt.Printf("gostatic %s\n", goopt.Version)
		return
	}

	if len(goopt.Args) == 0 {
		println(goopt.Usage())
		return
	}

	config := RetrieveGlobalConfig(goopt.Args[0])

	site := NewSite(config)
	site.Summary()
}