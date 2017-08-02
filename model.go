// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

type Model struct {
	Personas       map[string]Persona
	SystemContexts map[string]SystemContext
	Systems        map[string]System
	Containers     map[string]Container
	Components     map[string]Component
	Relationships  []Relationship
	Errors         []error
}

func newModel() *Model {
	m := new(Model)
	m.Personas = make(map[string]Persona)
	m.SystemContexts = make(map[string]SystemContext)
	m.Systems = make(map[string]System)
	m.Containers = make(map[string]Container)
	m.Components = make(map[string]Component)
	m.Relationships = make([]Relationship, 0)
	m.Errors = make([]error, 0)
	return m
}

type Persona struct {
	Name        string
	Description string
	Tags        []string
}

type SystemContext struct {
	Name            string
	Description     string
	CoreSystems     []string
	ExternalSystems []string
}

type System struct {
	Name        string
	Description string
	Tags        []string
}

type Container struct {
	System      string
	Name        string
	Description string
	Technology  string
	Tags        []string
}

type Component struct {
	Container   string
	Name        string
	Description string
	Technology  string
	Tags        []string
}

type Relationship struct {
	Source      string
	Description string
	Technology  string
	Destination string
	Tags        []string
}

func (m Model) SourceRelationships(source string) []Relationship {
	rels := make([]Relationship, 0)
	for _, r := range m.Relationships {
		if r.Source != source {
			continue
		}
		rels = append(rels, r)
	}
	return rels
}
