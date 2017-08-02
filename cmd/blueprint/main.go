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
		view = blueprint.NewGenericSystemContextView(model)
	} else {
		name := r.URL.Path[1:]
		switch r.URL.Query().Get("view") {
		case "context":
			sysCtx, ok := model.SystemContexts[name]
			if !ok {
				http.Error(w, "Model not found.", http.StatusNotFound)
				return
			}
			view = blueprint.NewSystemContextView(model, sysCtx)
		case "container":
			sys, ok := model.Systems[name]
			if !ok {
				http.Error(w, "Model not found.", http.StatusNotFound)
				return
			}
			view = blueprint.NewContainerView(model, sys)
		case "component":
			component, ok := model.Components[name]
			if !ok {
				http.Error(w, "Model not found.", http.StatusNotFound)
				return
			}
			view = blueprint.NewComponentView(model, component)
		default:
			http.Error(w, "Unknown view kind", http.StatusBadRequest)
			return
		}
	}

	err = blueprint.RenderPage(w, view, model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
