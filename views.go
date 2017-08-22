// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

import (
	"io"
)

// A View represents a specific subset of entities of a complete model.
// It can be rendered using RenderHTMLPage.
type View interface {
	Title() string
	Description() string
	dot(w io.Writer, model Model) error
}

type systemContextView struct {
	title           string
	description     string
	CoreSystems     []string
	ExternalSystems []string
	Personas        []string
}

func (v systemContextView) Description() string {
	return v.description
}

func (v systemContextView) Title() string {
	return "[System Context] " + v.title
}

type containerView struct {
	title       string
	description string
	System      string
	Containers  []string
	Systems     []string
	Personas    []string
}

func (v containerView) Description() string {
	return v.description
}

func (v containerView) Title() string {
	return "[Containers] " + v.title
}

type componentView struct {
	title       string
	description string
	Container   string
	Components  []string
	Containers  []string
	Systems     []string
}

func (v componentView) Description() string {
	return v.description
}

func (v componentView) Title() string {
	return "[Components] " + v.title
}

func (m Model) NewSystemContextView(sysCtx SystemContext) View {
	personas := make([]string, 0)
	for k := range m.Personas {
		personas = append(personas, k)
	}

	return systemContextView{
		title:           sysCtx.Name,
		description:     sysCtx.Description,
		CoreSystems:     sysCtx.CoreSystems,
		ExternalSystems: sysCtx.ExternalSystems,
		Personas:        personas,
	}

}

func (m Model) NewContainerView(sys System) View {
	containers := make([]string, 0)
	systems := make([]string, 0)
	personas := make([]string, 0)

	for k, c := range m.Containers {
		if c.System != sys.Name {
			continue
		}
		containers = append(containers, k)
		for _, r := range m.Relationships {
			if r.Source == c.Name {
				if _, ok := m.Systems[r.Destination]; ok {
					systems = append(systems, r.Destination)
				}
				if _, ok := m.Personas[r.Destination]; ok {
					personas = append(personas, r.Destination)
				}
			} else if r.Destination == c.Name {
				if _, ok := m.Systems[r.Source]; ok {
					systems = append(systems, r.Source)
				}
				if _, ok := m.Personas[r.Source]; ok {
					personas = append(personas, r.Source)
				}

			}
		}
	}

	return containerView{
		title:       sys.Name,
		description: sys.Description,
		System:      sys.Name,
		Containers:  containers,
		Systems:     systems,
		Personas:    personas,
	}
}

func (m Model) NewComponentView(cont Container) View {
	components := make([]string, 0)
	containers := make([]string, 0)
	systems := make([]string, 0)

	for k, c := range m.Components {
		if c.Container != cont.Name {
			continue
		}
		components = append(components, k)
		for _, r := range m.Relationships {
			if r.Source == c.Name {
				if _, ok := m.Systems[r.Destination]; ok {
					systems = append(systems, r.Destination)
				}
				if _, ok := m.Containers[r.Destination]; ok {
					containers = append(containers, r.Destination)
				}
			} else if r.Destination == c.Name {
				if _, ok := m.Systems[r.Source]; ok {
					systems = append(systems, r.Source)
				}
				if _, ok := m.Containers[r.Source]; ok {
					containers = append(containers, r.Source)
				}
			}
		}
	}

	return componentView{
		title:       cont.Name,
		description: cont.Description,
		Container:   cont.Name,
		Containers:  containers,
		Components:  components,
		Systems:     systems,
	}
}

func (m Model) NewGenericSystemContextView() View {
	systems := make([]string, 0)
	for k := range m.Systems {
		systems = append(systems, k)
	}

	personas := make([]string, 0)
	for k := range m.Personas {
		personas = append(personas, k)
	}

	return systemContextView{
		title:           "System Context Diagram",
		description:     "The complete system context diagram, containing all systems of the current project.",
		CoreSystems:     systems,
		ExternalSystems: []string{},
		Personas:        personas,
	}
}
