// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/pkg/browser"
	"github.com/urld/blueprint"
)

var project struct {
	path string
}

func main() {
	addr := flag.String("http", ":8080", "HTTP Service address")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("exactly 1 project path required")
		os.Exit(2)
	}
	project.path = flag.Arg(0)

	if strings.HasPrefix(*addr, ":") {
		browser.OpenURL("http://localhost" + *addr)

	} else {
		browser.OpenURL("http://" + *addr)
	}
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	model, err := blueprint.Parse(project.path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var view blueprint.View
	if r.URL.Path[1:] == "" {
		view = model.NewGenericSystemContextView()
	} else {
		viewKind, name := path.Split(r.URL.Path[1:])
		name = strings.TrimSuffix(name, ".html")
		switch path.Dir(viewKind) {
		case "contexts":
			sysCtx, ok := model.SystemContexts[name]
			if !ok {
				http.Error(w, "Model not found.", http.StatusNotFound)
				return
			}
			view = model.NewSystemContextView(sysCtx)
		case "containers":
			sys, ok := model.Systems[name]
			if !ok {
				http.Error(w, "Model not found.", http.StatusNotFound)
				return
			}
			view = model.NewContainerView(sys)
		case "components":
			container, ok := model.Containers[name]
			if !ok {
				http.Error(w, "Model not found.", http.StatusNotFound)
				return
			}
			view = model.NewComponentView(container)
		default:
			http.Error(w, "Unknown view kind: "+viewKind, http.StatusBadRequest)
			return
		}
	}

	err = blueprint.RenderHTMLPage(w, view, model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
