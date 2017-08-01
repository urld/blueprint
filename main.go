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
		browser.OpenURL("http://localhost" + *addr + "/systemcontext")

	} else {
		browser.OpenURL("http://" + *addr + "/systemcontext")
	}
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	//kind := r.URL.Path[1]
	//name := r.URL.Path[2]

	model, err := Parse(project.path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view := systemContextView(model)

	err = renderPage(w, view, model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func systemContextView(m Model) View {
	systems := make([]string, 0)
	for k := range m.systems {
		systems = append(systems, k)
	}

	personas := make([]string, 0)
	for k := range m.personas {
		personas = append(personas, k)
	}

	return SystemContextView{title: "System Context Diagram", description: "The complete system context diagram, containing all systems of the current project.", Systems: systems, Personas: personas}
}
