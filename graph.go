// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"text/template"
)

const dotTemplate = `
digraph "{{.Title}}" {
	ranksep="1.5 equally";
	nodesep="1.5 equally";
	node[fontcolor="white" fontsize=11 fontname="Sans" shape="box" style="filled,rounded" margin="0.20,0.20"];
	edge[fontcolor="dimgrey" color="dimgrey" fontsize=11 fontname="Sans"];

	subgraph cluster_core {
		color="#7b7b7b";
		style="dashed,rounded,bold";
		"__core__" [label="" style="invis" fixedsize="true" height="0.01" width="0.01"];
		{{- range .CoreNodes}}
		"{{.Name}}" [{{range $k, $v := .Attrs}} {{$k}}=<{{$v}}>{{end}} ];
		{{- end}}
	}

	// Fake edges to ensure cluster ranks
	"__top__" -> "__core__" -> "__bottom__" [style="invis"];
	{{- range .CoreNodes}}
	"__top__" -> "{{.Name}}" -> "__bottom__" [style="invis"];
	{{- end}}

	subgraph cluster_top {
		rank="min,same";
		style="invis";
		"__top__" [label="" style="invis" fixedsize="true" height="0.01" width="0.01"];
		{{- range .TopNodes}}
		"{{.Name}}" [{{range $k, $v := .Attrs}} {{$k}}=<{{$v}}>{{end}} ];
		{{- end}}
	}

	subgraph cluster_bottom {
		rank="max,same";
		style="invis";
		"__bottom__" [label="" style="invis" fixedsize="true" height="0.01" width="0.01"];
		{{- range .BottomNodes}}
		"{{.Name}}" [{{range $k, $v := .Attrs}} {{$k}}=<{{$v}}>{{end}} ];
		{{- end}}
	}

	// relationships
	{{- range .Edges}}
	"{{.Source}}" -> "{{.Destination}}" [{{range $k, $v := .Attrs}} {{$k}}=<{{$v}}>{{end}} ];
	{{- end}}
}
`

type graph struct {
	Title       string
	CoreNodes   []node
	TopNodes    []node
	BottomNodes []node
	Edges       []edge
}

type node struct {
	Name  string
	Attrs map[string]string
}

type edge struct {
	Source      string
	Destination string
	Attrs       map[string]string
}

func renderGraph(w io.Writer, view View, model Model) error {
	cmd := exec.Command("dot", "-Tsvg")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	out, err := cmd.StdoutPipe()
	errPipe, err := cmd.StderrPipe()
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = view.dot(in, model)
	//err = view.dot(io.MultiWriter(in, os.Stdout), model)
	if err != nil {
		return err
	}
	_ = in.Close()

	_, err = io.Copy(w, out)
	if err != nil {
		return err
	}

	errBuf := new(bytes.Buffer)
	_, err = io.Copy(errBuf, errPipe)
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf(strings.Join([]string{err.Error(), errBuf.String()}, "\n"))
	}
	if errBuf.Len() > 0 {
		return fmt.Errorf(errBuf.String())
	}
	return nil
}

func (v componentView) dot(w io.Writer, model Model) error {
	coreNodes := make([]node, 0)
	topNodes := make([]node, 0)
	bottomNodes := make([]node, 0)
	edges := make([]edge, 0)
	names := make([]string, 0)

	for _, name := range v.Components {
		cont, ok := model.Components[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		coreNodes = append(coreNodes, componentNode(cont))
		names = append(names, name)
	}

	for _, name := range v.Containers {
		cont, ok := model.Containers[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		topNodes = append(topNodes, containerNode(cont))
		names = append(names, name)
	}

	for _, name := range v.Systems {
		sys, ok := model.Systems[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		bottomNodes = append(bottomNodes, systemNode(sys))
		names = append(names, name)
	}

	edges = append(edges, relationshipEdges(model.FindRelationships(names)...)...)

	return genDot(w, graph{Title: v.title, CoreNodes: coreNodes, TopNodes: topNodes, BottomNodes: bottomNodes, Edges: edges})
}

func (v containerView) dot(w io.Writer, model Model) error {
	coreNodes := make([]node, 0)
	topNodes := make([]node, 0)
	bottomNodes := make([]node, 0)
	edges := make([]edge, 0)
	names := make([]string, 0)

	for _, name := range v.Containers {
		cont, ok := model.Containers[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		coreNodes = append(coreNodes, containerNode(cont))
		names = append(names, name)
	}

	for _, name := range v.Systems {
		sys, ok := model.Systems[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		bottomNodes = append(bottomNodes, systemNode(sys))
		names = append(names, name)
	}

	for _, name := range v.Personas {
		pers, ok := model.Personas[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		topNodes = append(topNodes, personaNode(pers))
		names = append(names, name)
	}

	edges = append(edges, relationshipEdges(model.FindRelationships(names)...)...)

	return genDot(w, graph{Title: v.title, CoreNodes: coreNodes, TopNodes: topNodes, BottomNodes: bottomNodes, Edges: edges})
}

func (v systemContextView) dot(w io.Writer, model Model) error {
	coreNodes := make([]node, 0)
	topNodes := make([]node, 0)
	bottomNodes := make([]node, 0)
	edges := make([]edge, 0)
	names := make([]string, 0)

	for _, name := range v.CoreSystems {
		sys, ok := model.Systems[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		coreNodes = append(coreNodes, systemNode(sys))
		names = append(names, name)
	}

	for _, name := range v.ExternalSystems {
		sys, ok := model.Systems[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		bottomNodes = append(bottomNodes, systemNode(sys))
		names = append(names, name)
	}

	for _, name := range v.Personas {
		pers, ok := model.Personas[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		topNodes = append(topNodes, personaNode(pers))
		names = append(names, name)
	}

	edges = append(edges, relationshipEdges(model.FindRelationships(names)...)...)

	return genDot(w, graph{Title: v.title, CoreNodes: coreNodes, TopNodes: topNodes, BottomNodes: bottomNodes, Edges: edges})
}

func genDot(w io.Writer, g graph) error {
	t := template.Must(template.New("dotTemplate").Parse(dotTemplate))
	return t.Execute(w, g)
}
