package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Konstantin8105/errors"
)

type Names struct {
	Package   string   // for example : "main"
	Type      string   // for example : "[]float64"
	CodeNew   string   // for example : "make([]float64, size)"
	CacheName string   // for example : "FloatsCache"
	Imports   []string // for example :
}

func main() {
	var (
		pkg  = flag.String("pkg", "main", "package name")
		typ  = flag.String("type", "[]float64", "type of generated object")
		cnew = flag.String("new", "make([]float64, size)", "create a new object")
		name = flag.String("name", "FloatsCache", "name of cache")
		imp  = flag.String("imports", "", "additional imports separated by comma.\n"+
			"For example: \"mat , sparse\"")
		help = flag.Bool("h", false, "print help information")
	)

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	// create a struct
	n := Names{
		Package:   *pkg,
		Type:      *typ,
		CodeNew:   *cnew,
		CacheName: *name,
		Imports:   strings.Split(*imp, ","),
	}

	// check input data
	{
		et := errors.New("Check input data")

		// clean spaces
		n.Package = strings.TrimSpace(n.Package)
		n.Type = strings.TrimSpace(n.Type)
		n.CodeNew = strings.TrimSpace(n.CodeNew)
		n.CacheName = strings.TrimSpace(n.CacheName)
		for i := range n.Imports {
			n.Imports[i] = strings.TrimSpace(n.Imports[i])
		}

		if len(n.Package) == 0 {
			et.Add(fmt.Errorf("Package name is empty"))
		}
		if len(n.Type) == 0 {
			et.Add(fmt.Errorf("Type is empty"))
		}
		if len(n.CodeNew) == 0 {
			et.Add(fmt.Errorf("Add generation new code"))
		}
		if len(n.CacheName) == 0 {
			et.Add(fmt.Errorf("Name of cache is empty"))
		}

		if et.IsError() {
			fmt.Fprintf(os.Stderr, "%v", et)
			os.Exit(-1)
		}
	}

	// remove empty imports
	{
	again:
		for i := range n.Imports {
			if len(n.Imports[i]) == 0 {
				n.Imports = append(n.Imports[:i], n.Imports[i+1:]...)
				goto again
			}
		}
	}

	// generate code to stdout
	t := template.Must(template.New("template").Parse(code))

	err := t.Execute(os.Stdout, n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing template: %v", err)
		os.Exit(-2)
	}
}
