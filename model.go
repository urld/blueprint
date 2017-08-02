// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package main

type Model struct {
	personas       map[string]Persona
	systemContexts map[string]SystemContext
	systems        map[string]System
	containers     map[string]Container
	components     map[string]Component
	relationships  []Relationship
	views          map[string]View
	errors         []error
}

func newModel() *Model {
	m := new(Model)
	m.personas = make(map[string]Persona)
	m.systemContexts = make(map[string]SystemContext)
	m.systems = make(map[string]System)
	m.containers = make(map[string]Container)
	m.components = make(map[string]Component)
	m.relationships = make([]Relationship, 0)
	m.errors = make([]error, 0)
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

func (m Model) Relationships(source string) []Relationship {
	rels := make([]Relationship, 0)
	for _, r := range m.relationships {
		if r.Source != source {
			continue
		}
		rels = append(rels, r)
	}
	return rels
}
