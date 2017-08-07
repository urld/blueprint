// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

// Model is the C4 architecture model representation of a project.
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

// A Persona that interacts with other entities of the software system.
type Persona struct {
	Name        string
	Description string
	Tags        []string
}

// A SystemContext defines a subset of Systems of the whole project that
// interact with each other.
type SystemContext struct {
	Name            string
	Description     string
	CoreSystems     []string
	ExternalSystems []string
}

// A System according to the C4 software architecture model.
type System struct {
	Name        string
	Description string
	Tags        []string
}

// A Container according to the C4 software architecture model.
type Container struct {
	System      string
	Name        string
	Description string
	Technology  string
	Tags        []string
}

// A Component according to the C4 software architecture model.
type Component struct {
	Container   string
	Name        string
	Description string
	Technology  string
	Tags        []string
}

// A Relationship between two arbitrary entities of the C4 model.
type Relationship struct {
	Source      string
	Description string
	Technology  string
	Destination string
	Tags        []string
}

// FindRelationships searches for relationships which are relevant for a given
// set of set of nodes. A relationship is considered relevant if both its
// Source and Destination are part of the given node set.
func (m Model) FindRelationships(nodes []string) []Relationship {
	nodeSet := make(map[string]bool)
	for _, n := range nodes {
		nodeSet[n] = true
	}

	rels := make([]Relationship, 0)
	for _, r := range m.Relationships {
		if nodeSet[r.Source] && nodeSet[r.Destination] {
			rels = append(rels, r)
		}
	}
	return rels
}
