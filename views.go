// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

import (
	"io"
)

type View interface {
	Title() string
	Description() string
	Dot(w io.Writer, model Model) error
}

type SystemContextView struct {
	title           string
	description     string
	CoreSystems     []string
	ExternalSystems []string
	Personas        []string
}

func (v SystemContextView) Description() string {
	return v.description
}

func (v SystemContextView) Title() string {
	return v.title
}

type ContainerView struct {
	title       string
	description string
	System      string
	Containers  []string
	Systems     []string
	Personas    []string
}

func (v ContainerView) Description() string {
	return v.description
}

func (v ContainerView) Title() string {
	return v.title
}

type ComponentView struct {
	title       string
	description string
	Container   string
	Components  []string
	Containers  []string
	Systems     []string
	Personas    []string
}

func (v ComponentView) Description() string {
	return v.description
}

func (v ComponentView) Title() string {
	return v.title
}

func NewSystemContextView(m Model, sysCtx SystemContext) View {
	personas := make([]string, 0)
	for k := range m.Personas {
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

func NewContainerView(m Model, sys System) View {
	containers := make([]string, 0)
	systems := make([]string, 0)
	for k, c := range m.Containers {
		if c.System != sys.Name {
			continue
		}
		containers = append(containers, k)
		for _, r := range m.SourceRelationships(c.Name) {
			if _, ok := m.Systems[r.Destination]; ok {
				systems = append(systems, r.Destination)
			}
		}
	}

	personas := make([]string, 0)
	for k := range m.Personas {
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

func NewComponentView(m Model, c Component) View {
	// TODO
	return nil
}

func NewGenericSystemContextView(m Model) View {
	systems := make([]string, 0)
	for k := range m.Systems {
		systems = append(systems, k)
	}

	personas := make([]string, 0)
	for k := range m.Personas {
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
