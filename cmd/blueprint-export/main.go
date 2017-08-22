// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/urld/blueprint"
)

var (
	projPath   string
	outputPath string
)

func main() {
	flag.StringVar(&projPath, "project", "", "path to project directory")
	flag.StringVar(&outputPath, "output", "", "path to output directory")
	flag.Parse()

	if projPath == "" || outputPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	err := renderProject()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func renderProject() error {
	model, err := blueprint.Parse(projPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(outputPath, "components"), 0755)
	if err != nil {
		return err
	}

	for name, container := range model.Containers {
		view := model.NewComponentView(container)

		err := write(path.Join(outputPath, "components", name+".html"), view, model)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(path.Join(outputPath, "containers"), 0755)
	if err != nil {
		return err
	}

	for name, system := range model.Systems {
		view := model.NewContainerView(system)

		err := write(path.Join(outputPath, "containers", name+".html"), view, model)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(path.Join(outputPath, "contexts"), 0755)
	if err != nil {
		return err
	}

	for name, context := range model.SystemContexts {
		view := model.NewSystemContextView(context)

		err := write(path.Join(outputPath, "contexts", name+".html"), view, model)
		if err != nil {
			return err
		}
	}

	view := model.NewGenericSystemContextView()
	err = write(path.Join(outputPath, "contexts", "index.html"), view, model)
	if err != nil {
		return err
	}

	return nil
}

func write(filePath string, view blueprint.View, model blueprint.Model) error {
	f, err := os.Create(filePath)
	defer close(f)
	if err != nil {
		return err
	}

	err = blueprint.RenderHTMLPage(f, view, model)
	if err != nil {
		return err
	}

	return nil
}

func close(c io.Closer) {
	_ = c.Close()
}
