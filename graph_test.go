// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

import (
	"bytes"
	"testing"
)

const expectedGenDot = `
digraph "Test Title" {
	ranksep="1.5 equally";
	nodesep="1.5 equally";
	node[fontcolor="white" fontsize=11 fontname="Sans" shape="box" style="filled,rounded" margin="0.20,0.20"];
	edge[fontcolor="dimgrey" color="dimgrey" fontsize=11 fontname="Sans"];

	subgraph cluster_core {
		color="#7b7b7b";
		style="dashed,rounded,bold";
		"N1" [ label=<N1 Label> style=<filled> ];
		"N2" [ ];
	}

	// other nodes:
	"N3" [ label=<N3 Label> ];

	// relationships
	"N1" -> "N2" [ label=<1-2 Label> ];
}
`

func TestGenDot(t *testing.T) {

	n := []node{{Name: "N1", Attrs: map[string]string{"label": "N1 Label", "style": "filled"}},
		{Name: "N2", Attrs: map[string]string{}}}
	en := []node{{Name: "N3", Attrs: map[string]string{"label": "N3 Label"}}}
	e := []edge{{Source: "N1", Destination: "N2", Attrs: map[string]string{"label": "1-2 Label"}}}
	g := graph{Title: "Test Title", CoreNodes: n, ExternalNodes: en, Edges: e}

	buf := new(bytes.Buffer)
	genDot(buf, g)

	assertEqual(t, expectedGenDot, buf.String(), "generated dot input does not match")
}
