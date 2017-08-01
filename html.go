// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package main

import (
	"bytes"
	"html/template"
	"io"
)

const pageTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>{{.Title}}</title>
	<style>
	body {
		font-family: Sans;
	}
	.panel-margin {
		margin-top: 16px;
		margin-bottom: 16px;
		margin-left: 16px;
		margin-right: 22px;
	}
	.panel {
		margin-top: 16px;
		margin-bottom: 16px;
		overflow: hidden
	}
	.danger {
		background-color: #ffdddd;
		border-left: 6px solid #f44336;
	}
	.success {
		background-color: #ddffdd;
		border-left: 6px solid #4CAF50;
	}
	.info {
		background-color: #e7f3fe;
		border-left: 6px solid #2196F3;
	}
	.warning {
		background-color: #ffffcc;
		border-left: 6px solid #ffeb3b;
	}
	</style>
</head>
<body>
	<h1>{{.Title}}</h1>
	<p>{{.Description}}</p>

	{{if .ModelErrors -}}
	<div><div class="danger panel"><div class="panel-margin">
		Model Errors
		<pre>{{range .ModelErrors}}{{.Error}}<br/>{{end}}</pre>
	</div></div></div>
	{{- end}}

	{{if .GenError -}}
	<div><div class="danger panel"><div class="panel-margin">
		Graphviz Errors
		<pre>{{.GenError.Error}}</pre>
	</div></div></div>
	{{- end}}

	{{if .GenWarning -}}
	<div><div class="warning panel"><div class="panel-margin">
		Graphviz Warnings
		<pre>{{.GenWarning.Error}}</pre>
	</div></div></div>
	{{- end}}

	<div>{{.Svg}}</div>
</body>
</html>
`

type page struct {
	Title       string
	Description string
	Svg         template.HTML
	ModelErrors []error
	GenError    error
	GenWarning  error
}

func renderPage(w io.Writer, view View, model Model) error {
	p := page{
		Title:       view.Title(),
		Description: view.Description(),
		ModelErrors: model.errors,
	}

	svgBuf := new(bytes.Buffer)
	err := renderGraph(svgBuf, view, model)
	if err != nil && svgBuf.Len() == 0 {
		p.GenError = err
	} else if err != nil {
		p.GenWarning = err
	}

	p.Svg = template.HTML(svgBuf.String())

	t := template.Must(template.New("pageTemplate").Parse(pageTemplate))
	return t.Execute(w, p)
}
