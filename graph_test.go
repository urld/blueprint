// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"
)

const expectedGenDot = `
digraph "Test Title" {
	ranksep="1.5 equally";
	splines="polyline";
	node[fontcolor="white" fontsize=11 fontname="Sans" shape="box" style="filled,rounded" margin="0.20,0.20"];
	edge[fontcolor="dimgrey" color="dimgrey" fontsize=11 fontname="Sans"];

	"N1" [ label=<N1 Label> style=<filled> ];
	"N2" [ ];
	
	"N1" -> "N2" [ label=<1-2 Label> ];
	
}
`

func TestGenDot(t *testing.T) {

	n := []node{node{Name: "N1", Attrs: map[string]string{"label": "N1 Label", "style": "filled"}},
		node{Name: "N2", Attrs: map[string]string{}}}
	e := []edge{edge{Source: "N1", Destination: "N2", Attrs: map[string]string{"label": "1-2 Label"}}}
	g := graph{Title: "Test Title", Nodes: n, Edges: e}

	buf := new(bytes.Buffer)
	genDot(buf, g)

	assertEqual(t, expectedGenDot, buf.String(), "generated dot input does not match")
}
