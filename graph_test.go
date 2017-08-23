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
	ranksep="1.0";
	nodesep="1.0";
	node[fontcolor="white" fontsize=11 fontname="Sans" shape="box" style="filled,rounded" margin="0.20,0.20"];
	edge[fontcolor="dimgrey" color="dimgrey" fontsize=11 fontname="Sans"];

	subgraph cluster_core {
		color="#7b7b7b";
		style="dashed,rounded,bold";
		"N1" [ label=<N1 Label> style=<filled> ];
		"N2" [ ];
	}

	subgraph cluster_top {
		rank="sink";
		style="invis";
		"N3" [ label=<N3 Label> ];
	}

	subgraph cluster_bottom {
		rank="source";
		style="invis";
		"N4" [ label=<N4 Label> ];
	}

	// relationships
	"N1" -> "N2" [ label=<1-2 Label> ];
}
`

func TestGenDot(t *testing.T) {

	n := []node{{Name: "N1", Attrs: map[string]string{"label": "N1 Label", "style": "filled"}},
		{Name: "N2", Attrs: map[string]string{}}}
	tn := []node{{Name: "N3", Attrs: map[string]string{"label": "N3 Label"}}}
	bn := []node{{Name: "N4", Attrs: map[string]string{"label": "N4 Label"}}}
	e := []edge{{Source: "N1", Destination: "N2", Attrs: map[string]string{"label": "1-2 Label"}}}
	g := graph{Title: "Test Title", CoreNodes: n, TopNodes: tn, BottomNodes: bn, Edges: e}

	buf := new(bytes.Buffer)
	err := genDot(buf, g)

	assertEqual(t, expectedGenDot, buf.String(), "generated dot input does not match")
	assertEqual(t, nil, err, "genDot returned an error")
}
