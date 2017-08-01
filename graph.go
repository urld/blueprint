// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package main

import (
	"io"
	"os"
	"os/exec"
	"text/template"
)

const (
	systemColor          = "#08427b"
	systemBorderColor    = "#002a56"
	personColor          = "#08427b"
	personBorderColor    = "#002a56"
	containerColor       = "#1168bd"
	containerBorderColor = "#065fae"
	componentColor       = "#3d88d1"
	componentBorderColor = "#1782cc"

	lineSep   = "<BR/>"
	lineLimit = 38
)

const dotTemplate = `
digraph "{{.Title}}" {
	ranksep="1.5 equally";
	splines="polyline";
	node[fontcolor="white" fontsize=11 fontname="Sans" shape="box" style="filled,rounded" margin="0.20,0.20"];
	edge[fontcolor="dimgrey" color="dimgrey" fontsize=11 fontname="Sans"];

	{{range .Nodes -}}
		"{{.Name}}" [{{range $k, $v := .Attrs}} {{$k}}=<{{$v}}>{{end}} ];
	{{end}}
	{{range .Edges -}}
		"{{.Source}}" -> "{{.Destination}}" [{{range $k, $v := .Attrs}} {{$k}}=<{{$v}}>{{end}} ];
	{{end}}
}
`

type graph struct {
	Title string
	Nodes []node
	Edges []edge
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
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = view.Dot(io.MultiWriter(in, os.Stdout), model)
	if err != nil {
		return err
	}
	_ = in.Close()

	_, err = io.Copy(w, out)
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func systemNode(s System) node {
	attrs := map[string]string{
		"label": "<FONT POINT-SIZE=\"14\"><B>" + s.Name + "</B></FONT><BR/>[System]<BR/><BR/>" +
			WrapWords(s.Description, lineSep, lineLimit),
		"fillcolor": systemColor,
		"color":     systemBorderColor,
	}
	return node{Name: s.Name, Attrs: attrs}
}

func personaNode(p Persona) node {
	attrs := map[string]string{
		"label": "<FONT POINT-SIZE=\"14\"><B>" + p.Name + "</B></FONT><BR/>[Person]<BR/><BR/>" +
			WrapWords(p.Description, lineSep, lineLimit),
		"fillcolor": personColor,
		"color":     personBorderColor,
	}
	return node{Name: p.Name, Attrs: attrs}
}

func relationshipEdge(r Relationship) edge {
	attrs := map[string]string{
		"label": WrapWords(r.Description, lineSep, lineLimit),
	}
	return edge{Source: r.Source, Destination: r.Destination, Attrs: attrs}
}

func relationshipEdges(rs ...Relationship) []edge {
	edges := make([]edge, 0)
	for _, r := range rs {
		edges = append(edges, relationshipEdge(r))
	}
	return edges
}

func (v SystemContextView) Dot(w io.Writer, model Model) error {
	nodes := make([]node, 0)
	edges := make([]edge, 0)

	for _, name := range v.Systems {
		sys, ok := model.systems[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		nodes = append(nodes, systemNode(sys))
		edges = append(edges, relationshipEdges(model.Relationships(name)...)...)
	}

	for _, name := range v.Personas {
		pers, ok := model.personas[name]
		if !ok {
			// TODO: is this an error?
			continue
		}
		nodes = append(nodes, personaNode(pers))
		edges = append(edges, relationshipEdges(model.Relationships(name)...)...)
	}

	return genDot(w, graph{Title: v.title, Nodes: nodes, Edges: edges})
}

func genDot(w io.Writer, g graph) error {
	t := template.Must(template.New("dotTemplate").Parse(dotTemplate))
	return t.Execute(w, g)
}
