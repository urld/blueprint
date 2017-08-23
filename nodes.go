// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

import (
	"net/url"
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

	lineLimit = 38
)

func systemNode(s System) node {
	attrs := map[string]string{
		"label": "<FONT POINT-SIZE=\"14\"><B>" + s.Name + "</B></FONT><BR/>" +
			"[System]<BR/><BR/>" +
			wrapWords(s.Description, lineLimit),
		"fillcolor": systemColor,
		"color":     systemBorderColor,
		"URL":       "../containers/" + url.PathEscape(s.Name) + ".html",
	}
	return node{Name: s.Name, Attrs: attrs}
}

func containerNode(c Container) node {
	attrs := map[string]string{
		"label": "<FONT POINT-SIZE=\"14\"><B>" + c.Name + "</B></FONT><BR/>" +
			nodeTechnology("Container", c.Technology) + "<BR/><BR/>" +
			wrapWords(c.Description, lineLimit),
		"fillcolor": containerColor,
		"color":     containerBorderColor,
		"URL":       "../components/" + url.PathEscape(c.Name) + ".html",
	}
	return node{Name: c.Name, Attrs: attrs}
}

func componentNode(c Component) node {
	attrs := map[string]string{
		"label": "<FONT POINT-SIZE=\"14\"><B>" + c.Name + "</B></FONT><BR/>" +
			nodeTechnology("Component", c.Technology) + "<BR/><BR/>" +
			wrapWords(c.Description, lineLimit),
		"fillcolor": componentColor,
		"color":     componentBorderColor,
	}
	return node{Name: c.Name, Attrs: attrs}
}

func personaNode(p Persona) node {
	attrs := map[string]string{
		"label": "<FONT POINT-SIZE=\"14\"><B>" + p.Name + "</B></FONT><BR/>" +
			"[Persona]<BR/><BR/>" +
			wrapWords(p.Description, lineLimit),
		"fillcolor": personColor,
		"color":     personBorderColor,
	}
	return node{Name: p.Name, Attrs: attrs}
}

func relationshipEdge(r Relationship) edge {
	attrs := map[string]string{
		"label": "<TABLE BORDER=\"0\"><TR><TD>" + wrapWords(r.Description, lineLimit) + edgeTechnology(r) + "</TD></TR></TABLE>",
	}
	return edge{Source: r.Source, Destination: r.Destination, Attrs: attrs}
}

func edgeTechnology(r Relationship) string {
	if r.Technology == "" {
		return ""
	}
	return "<BR/>[" + wrapWords(r.Technology, lineLimit) + "]"
}

func nodeTechnology(nodeKind, technology string) string {
	if technology == "" {
		return "[" + nodeKind + "]"
	}
	return "[" + wrapWords(nodeKind+": "+technology, lineLimit) + "]"
}

func relationshipEdges(rs ...Relationship) []edge {
	edges := make([]edge, 0)
	for _, r := range rs {
		edges = append(edges, relationshipEdge(r))
	}
	return edges
}
