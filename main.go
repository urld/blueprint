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
	model, err := Parse(project.path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var view View
	if r.URL.Path[1:] == "" {
		view = genericSystemContextView(model)
	} else {
		name := r.URL.Path[1:]
		switch r.URL.Query().Get("view") {
		case "context":
			sysCtx, ok := model.systemContexts[name]
			if !ok {
				http.Error(w, "Model not found.", http.StatusNotFound)
				return
			}
			view = systemContextView(model, sysCtx)
		case "container":
			sys, ok := model.systems[name]
			if !ok {
				http.Error(w, "Model not found.", http.StatusNotFound)
				return
			}
			view = containerView(model, sys)
		case "component":
			component, ok := model.components[name]
			if !ok {
				http.Error(w, "Model not found.", http.StatusNotFound)
				return
			}
			view = componentView(model, component)
		default:
			http.Error(w, "Unknown view kind", http.StatusBadRequest)
			return
		}
	}

	err = renderPage(w, view, model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func systemContextView(m Model, sysCtx SystemContext) View {
	personas := make([]string, 0)
	for k := range m.personas {
		personas = append(personas, k)
	}

	return SystemContextView{
		title:           sysCtx.Name,
		description:     sysCtx.Description,
		CoreSystems:     sysCtx.CoreSystems,
		ExternalSystems: sysCtx.ExternalSystems,
		Personas:        personas,
	}

}

func containerView(m Model, sys System) View {
	containers := make([]string, 0)
	systems := make([]string, 0)
	for k, c := range m.containers {
		if c.System != sys.Name {
			continue
		}
		containers = append(containers, k)
		for _, r := range m.Relationships(c.Name) {
			if _, ok := m.systems[r.Destination]; ok {
				systems = append(systems, r.Destination)
			}
		}
	}

	personas := make([]string, 0)
	for k := range m.personas {
		personas = append(personas, k)
	}

	return ContainerView{
		title:       sys.Name,
		description: sys.Description,
		System:      sys.Name,
		Containers:  containers,
		Systems:     systems,
		Personas:    personas,
	}
}

func componentView(m Model, c Component) View {
	// TODO
	return nil
}

func genericSystemContextView(m Model) View {
	systems := make([]string, 0)
	for k := range m.systems {
		systems = append(systems, k)
	}

	personas := make([]string, 0)
	for k := range m.personas {
		personas = append(personas, k)
	}

	return SystemContextView{
		title:           "System Context Diagram",
		description:     "The complete system context diagram, containing all systems of the current project.",
		CoreSystems:     systems,
		ExternalSystems: []string{},
		Personas:        personas,
	}
}
